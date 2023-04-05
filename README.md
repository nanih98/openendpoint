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


# Documentation
Take a look inside [docs](./docs) folder.

# Badges

![Build Status](https://github.com/nanih98/openendpoint/actions/workflows/releases.yml/badge.svg)
![Linter Status](https://github.com/nanih98/openendpoint/actions/workflows/lint.yml/badge.svg)
[![License](https://img.shields.io/github/license/nanih98/openendpoint)](/LICENSE)
[![Release](https://img.shields.io/github/release/nanih98/openendpoint)](https://github.com/nanih98/openendpoint/releases/latest)
[![GitHub Releases Stats](https://img.shields.io/github/downloads/nanih98/openendpoint/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=nanih98&repository=openendpoint)

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