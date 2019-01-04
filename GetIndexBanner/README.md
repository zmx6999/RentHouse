# GetIndexBanner Service

This is the GetIndexBanner service

Generated with

```
micro new sss/GetIndexBanner --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.GetIndexBanner
- Type: srv
- Alias: GetIndexBanner

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./GetIndexBanner-srv
```

Build a docker image
```
make docker
```