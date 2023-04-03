<p align="center" >
    <img src="assets/logo.png" alt="logo" width="250"/>
<h3 align="center">openendpoint (ALPHA)</h3>
<p align="center">Scan cloud public endpoints</p>
</p>

<p align="center" >
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/nanih98/openendpoint">
    <img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/nanih98/openendpoint">
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/nanih98/openendpoint">
</p>

<p align="center" >
    <a href="https://github.com/nanih98/openendpoint/actions/workflows/releases.yml"><img alt="Pipeline" src="https://github.com/nanih98/openendpoint/actions/workflows/releases.yml/badge.svg"></a>
    <a href="https://github.com/nanih98/openendpoint/actions/workflows/lint.yml"><img alt="Pipeline" src="https://github.com/nanih98/openendpoint/actions/workflows/lint.yml/badge.svg"></a>
    <a href="https://github.com/nanih98/openendpoint/blob/main/LICENSE"><img alt="LICENSE" src="https://img.shields.io/github/license/nanih98/openendpoint"></a>
</p>

# Usage
```sh
$ git clone https://github.com/nanih98/openendpoint.git
$ go run cmd/openendpoint/main.go -k test -w 10 -n 1.1.1.1 -f assets/fuzz.txt
```

# Credits

- [Folder layout](https://github.com/golang-standards/project-layout)
- [Inspired in cloud enum tool](https://github.com/initstring/cloud_enum)
- [Zap logging library](https://github.com/uber-go/zap)

# Contribution

Pull requests are welcome! Any code refactoring, improvement, implementation...

# LICENSE

[LICENSE](./LICENSE)