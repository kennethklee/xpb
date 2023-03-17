package xpb

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"
)

type module struct {
	Module      string
	Version     string
	Replacement string
}

/**
 * Returns the module string for display
 * i.e. github.com/pocketbase/pocketbase/core@master
 */
func (p module) String() string {
	return p.Module + "@" + p.Version
}

/**
 * Returns the trimmed module string for go get and go mod edit
 * i.e. github.com/pocketbase/pocketbase@master
 */
func (p module) AsModuleString() string {
	goModPathParts := strings.Split(p.Module, "/")
	goModPath := strings.Join(goModPathParts[:3], "/")
	return goModPath + "@" + p.Version
}

type Builder struct {
	TempProjectDir string // path to the temp custom pocketbase project directory

	pb      string   // pocketbase version
	modules []module // modules for plugins
}

func NewBuilder(pbVersion string, pluginStrs ...string) (*Builder, error) {
	plugins := []module{}
	for _, pluginStr := range pluginStrs {
		plugin, err := parsePlugin(pluginStr)
		if err != nil {
			return nil, err
		}

		plugins = append(plugins, plugin)
	}

	return &Builder{pb: pbVersion, modules: plugins}, nil
}

func parsePlugin(plugin string) (result module, err error) {
	// module[@version][=replacement]
	pluginPattern, err := regexp.Compile(`^(?P<module>[^=@]+)(@(?P<version>[^=]+))?(=(?P<replacement>.+))?$`)
	if err != nil {
		return
	}

	match := pluginPattern.FindStringSubmatch(plugin)
	if len(match) == 0 {
		err = fmt.Errorf("invalid plugin: %s", plugin)
		return
	}

	// Defaults
	result.Version = "latest"
	for i, name := range pluginPattern.SubexpNames() {
		if match[i] == "" {
			continue
		}
		switch name {
		case "module":
			result.Module = match[i]
		case "version":
			result.Version = match[i]
		case "replacement":
			result.Replacement = match[i]
		}
	}
	fmt.Println("[INFO]", "Plugin:", result.Module, result.Version, result.Replacement)
	return
}

func runGoMod(projectPath string, args ...string) error {
	args = append([]string{"mod"}, args...)
	goCmd := exec.Command("go", args...)
	goCmd.Dir = projectPath
	goCmd.Stdout = os.Stdout
	goCmd.Stderr = os.Stderr
	goCmd.Stdin = os.Stdin
	return goCmd.Run()
}

func runGoGet(projectPath string, args ...string) error {
	args = append([]string{"get"}, args...)
	goCmd := exec.Command("go", args...)
	goCmd.Dir = projectPath
	goCmd.Stdout = os.Stdout
	goCmd.Stderr = os.Stderr
	goCmd.Stdin = os.Stdin
	fmt.Println("[INFO]", goCmd.String())
	return goCmd.Run()
}

func (b *Builder) build() error {
	// Generate temp go project in /tmp
	projectPath, err := os.MkdirTemp("", "pocketbase")
	if err != nil {
		return fmt.Errorf("failed to create temp project path: %w", err)
	}
	// fmt.Println("[INFO]", "Temp project path:", projectPath)
	b.TempProjectDir = projectPath

	// Create project `go mod init`
	err = runGoMod(projectPath, "init", "pocketbase")
	if err != nil {
		return fmt.Errorf("failed to create temp project: %w", err)
	}

	// Go replace plugins `go mod edit -replace module@version=replacement`
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	for _, module := range b.modules {
		if module.Replacement != "" {
			replacement := path.Join(pwd, module.Replacement)
			err := runGoMod(projectPath, "edit", "-replace", module.AsModuleString()+"="+replacement)
			if err != nil {
				return fmt.Errorf("failed to replace plugin %s: %w", module.Module, err)
			}
		}
	}

	// Generate main.go (import main, plugins, and call main)
	tpl, err := template.New("main.go").Parse(mainGoTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse main.go template: %w", err)
	}

	mainGoFile, err := os.Create(projectPath + "/main.go")
	if err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	defer mainGoFile.Close()

	err = tpl.Execute(mainGoFile, struct{ Plugins []module }{Plugins: b.modules})
	if err != nil {
		return fmt.Errorf("failed to execute main.go template: %w", err)
	}

	// Go mod tidy `go mod tidy`
	err = runGoMod(projectPath, "tidy")
	if err != nil {
		return fmt.Errorf("failed to tidy project: %w", err)
	}

	// Go get pocketbase `go get github.com/pocketbase/pocketbase@<version>`
	err = runGoGet(projectPath, "github.com/pocketbase/pocketbase@"+b.pb)
	if err != nil {
		return fmt.Errorf("failed to get pocketbase@%s: %w", b.pb, err)
	}

	return nil
}

func (b *Builder) Compile(buildArgs ...string) error {
	err := b.build()
	if err != nil {
		return err
	}

	// Build `go build -o pocketbase main.go`
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	buildArgs = append(buildArgs, ".")
	buildArgs = append([]string{"build", "-o", path.Join(pwd, "pocketbase")}, buildArgs...)
	fmt.Println("[INFO] RUN:", "go", buildArgs)
	goBuildCmd := exec.Command("go", buildArgs...)
	goBuildCmd.Env = os.Environ()
	goBuildCmd.Dir = b.TempProjectDir
	goBuildCmd.Stdout = os.Stdout
	goBuildCmd.Stderr = os.Stderr
	return goBuildCmd.Run()
}

func (b *Builder) Run(args ...string) error {
	err := b.build()
	if err != nil {
		// Wait for input to see error
		fmt.Println("[ERROR]", "Failed to tidy project. Press enter to continue")
		fmt.Scanln()
		return err
	}

	// Run `go run main.go args`
	goRunCmd := exec.Command("go", "run", "main.go")
	goRunCmd.Env = os.Environ()
	goRunCmd.Dir = b.TempProjectDir
	goRunCmd.Args = append(goRunCmd.Args, args...)
	goRunCmd.Stdout = os.Stdout
	goRunCmd.Stderr = os.Stderr
	goRunCmd.Stdin = os.Stdin
	err = goRunCmd.Run()

	// Wait for input to see error
	if err != nil {
		fmt.Println("[ERROR]", "Failed to run project. Press enter to continue")
		fmt.Scanln()
	}
	return err
}

func (b *Builder) Close() error {
	if b.TempProjectDir != "" {
		return os.RemoveAll(b.TempProjectDir)
	}
	return nil
}

const mainGoTemplate = `package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/pocketbase/pocketbase"
	"github.com/kennethklee/xpb"

	{{ range .Plugins }}
	_ "{{.Module}}"
	{{ end }}
)

func main() {
	err := xpb.FireOnPreload()
	if err != nil {
		log.Fatal(err)
	}

	var app = pocketbase.New()

	plugins := xpb.GetPlugins()
	bold := color.New(color.Bold).Add(color.FgGreen)
	if len(plugins) > 0 {
		bold.Println("> Plugins")
		for _, plugin := range plugins {
			pluginInfo := plugin.Info()
			fmt.Printf("  - %s (%s) %s\n", pluginInfo.Name, pluginInfo.Version, pluginInfo.Description)
		}
	}

	err = xpb.FireOnLoad(app)
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Start(); err != nil {
		log.Fatal(err)
	}
}
`
