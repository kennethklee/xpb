# xpocketbase

Plugin orientated custom pocketbase builder, inspired by [xcaddy](https://github.com/caddyserver/xcaddy/).

This is a usable example of running pocketbase

## Requirements

- Golang

## Install

Install from source:

```
go install github.com/kennethklee/xpb/xpocketbase@latest
```

## Usage

`xpocketbase build <version>` will build a `pocketbase` binary.

```
# Build a specific version of pocketbase
xpocketbase build master
xpocketbase build v0.8.0

# Build the examples/base pocketbase
xpocketbase build latest \
    --with github.com/kennethklee/xpb/plugins/static \
    --with github.com/kennethklee/xpb/plugins/migrations \
    --with github.com/kennethklee/xpb/plugins/timeout

# Build with plugin module in current directory
xpocketbase build latest \
    --with my-module=.

# Specific version
xpocketbase build latest \
    --with github.com/kennethklee/xpb/plugins/static@v1.0.0

# Replaces contents of module with contents elsewhere
xpocketbase build latest \
    --with github.com/kennethklee/xpb/plugins/static@v1.0.0=../plugin-fork

# go build flags
xpocketbase build latest \
    --with github.com/kennethklee/xpb/plugins/static@v1.0.0 \
    -- -ldflags "-s -w"
```