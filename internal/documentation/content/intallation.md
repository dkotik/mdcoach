---
title: "Installation"
weight: 20
---

## Requirements

- Go 1.25 or later.
- A browser to view generated HTML presentations.

The compiler can fetch remote image URLs, so network access is needed when a source presentation uses remote images. Local images can be used without network access.

## Install the command-line tool

```sh
go install github.com/dkotik/mdcoach/cmd/mdcoach@latest
```

Make sure Go's binary directory is on your `PATH`, then verify the installation:

```sh
mdcoach --help
```

## Build from source

```sh
git clone https://github.com/dkotik/mdcoach.git
cd mdcoach
go build ./cmd/mdcoach
```

Run the tests with:

```sh
go test ./...
```

## Use in another Go module

Import the presentation package and add the module as a dependency:

```sh
go get github.com/dkotik/mdcoach/presentation
```

See [Using MdCoach as a Go library](/library.html) for a minimal example.
