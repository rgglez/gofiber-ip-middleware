# gofiber-ip-middleware

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub all releases](https://img.shields.io/github/downloads/rgglez/gofiber-ip-middleware/total)
![GitHub issues](https://img.shields.io/github/issues/rgglez/gofiber-ip-middleware)
![GitHub commit activity](https://img.shields.io/github/commit-activity/y/rgglez/gofiber-ip-middleware)
[![Go Report Card](https://goreportcard.com/badge/github.com/rgglez/gofiber-ip-middleware/src)](https://goreportcard.com/report/github.com/rgglez/gofiber-ip-middleware/src)
[![GitHub release](https://img.shields.io/github/release/rgglez/gofiber-ip-middleware.svg)](https://github.com/rgglez/gofiber-ip-middleware/releases/)
![GitHub stars](https://img.shields.io/github/stars/rgglez/gofiber-ip-middleware?style=social)
![GitHub forks](https://img.shields.io/github/forks/rgglez/gofiber-ip-middleware?style=social)

**gofiber-ip-middleware** is a [gofiber](https://gofiber.io/) [middleware](https://docs.gofiber.io/category/-middleware/) to verify if the client IP address 
(or the x-forwarded-for) is in a list of allowed addresses (or even CIDR blocks or hostnames).

## Installation

```bash
go get github.com/rgglez/gofiber-ip-middleware
```

## Usage

```go
import gofiberip "github.com/rgglez/gofiber-ip-middleware/gofiberip"

// Initialize Fiber app and middleware
app := fiber.New()
app.Use(gofiberip.New(gofiberzitadel.Config{AllowedIPs: []string{"267.132.21.1"}}))
```

## Configuration

There are some configuration options available in the ```Config``` struct:

* **```Next```** defines a function to skip this middleware when returned true. Optional. Default: nil
* **```AllowedIPs```** an array of strings, which could be IP addresses, CIDR blocks or hostnames. Required.

## Example

An example is included in the [example](example/) directory.


## Dependencies

* [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2)

## License

Copyright (c) 2025 Rodolfo González González

Licensed under the [Apache 2.0](LICENSE) license. Read the [LICENSE](LICENSE) file.

