# erry

[![golangci-lint](https://github.com/k1gabyt0/erry/actions/workflows/golangci-lint.yaml/badge.svg)](https://github.com/k1gabyt0/erry/actions/workflows/golangci-lint.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/k1gabyt0/erry.svg)](https://pkg.go.dev/github.com/k1gabyt0/erry)

## Description

A Go package for handling multiple errors.


⚡️ The project is in rapid stage of development, issues and critics are welcome.

## Features

- simple minded
- standart library errors.As & errors.Is compatible
- dependency free

## Table of contents

- [erry](#erry)
  - [Description](#description)
  - [Features](#features)
  - [Table of contents](#table-of-contents)
  - [Getting started](#getting-started)
    - [Installation](#installation)
    - [Usage](#usage)
      - [Create new multi-error](#create-new-multi-error)
      - [Transform existing error into multi-error](#transform-existing-error-into-multi-error)
  - [Alternatives](#alternatives)

## Getting started

See the [documentation](https://pkg.go.dev/github.com/k1gabyt0/erry) for further details.

### Installation

```bash
go get github.com/k1gabyt0/erry
```

### Usage

#### Create new multi-error

```go
errA := errors.New("err A")
errB := errors.New("err B")

multierr := erry.NewError("multierror", errA, errB)
fmt.Println(multierr)
```

**Output**:

```text
multierror:
    err A
    err B
This is error A
This is error B
```

#### Transform existing error into multi-error

```go
errA := errors.New("err A")
errB := errors.New("err B")

// Transforms errA into multi-error with errB
// as one of inner errors.
multierr := erry.ErrorFrom(errA, errB)
fmt.Println(multierr)

if errors.Is(multierr, errA) {
    fmt.Println("This is error A")
}
if errors.Is(multierr, errB) {
    fmt.Println("This is error B")
}
```

**Output**:

```text
err A:
    err B
This is error A
This is error B
```

## Alternatives

- [multierr](https://github.com/uber-go/multierr)
- [go-mutlierror](https://github.com/hashicorp/go-multierror)
- [talescale's multierr](https://github.com/tailscale/tailscale/tree/main/util/multierr)
