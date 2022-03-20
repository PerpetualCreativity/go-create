# go-create

An easy-to-use, concurrent cross-compilation wrapper around `go build`.

## Installation

```sh
go install github.com/PerpetualCreativity/go-create@latest
```

## Usage

[![asciicast](https://asciinema.org/a/pZEyWocq2RhoJeb44AE5ggRKl.svg)](https://asciinema.org/a/pZEyWocq2RhoJeb44AE5ggRKl)

### Explanation

```sh
go-create --name NAME --file main.go --compress --output=builds --exclude js android
```

This builds `main.go` for all OSs and architectures Go supports, except for JavaScript and Android. An example output filename is `./builds/NAME-linux-amd64`.

## Flags

| flag            | default                 | what it does                                  |
|-----------------|-------------------------|-----------------------------------------------|
| `--name`        | *none*, flag required   | specify name of builds                        |
| `--file`        | `main.go`               | specify file to build                         |
| `--compress`    | `false`                 | compress output                               |
| `--output`      | `.` (current directory) | directory to place builds                     |
| `--first-class` | `false`                 | include platforms with first-class Go support |
| `--cgo`         | `false`                 | include platforms with Cgo support            |
| `--exclude`     | `false`                 | see below                                     |


### `--exclude`

Normally, when running `go-create`, adding OSs at the end of the incantation will build for those OSs only. For example, if you only want to build for Linux and Darwin (macOS):

```sh
go-create --name [name here] linux darwin
```

But if you specify `--exclude`, you will get builds for every OS Go supports except for those at the end of the incantation. For example, if you want to build for everything except JavaScript (WebAssembly) and Android:

```sh
go-create --name [name here] --exclude js android
```

You can access a list of OSs and architectures that your Go installation supports using `go tool dist list`.

