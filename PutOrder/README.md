# PutOrder Service

This is the PutOrder service

Generated with

```
micro new sss/PutOrder --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.PutOrder
- Type: srv
- Alias: PutOrder

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
./PutOrder-srv
```

Build a docker image
```
make docker
```