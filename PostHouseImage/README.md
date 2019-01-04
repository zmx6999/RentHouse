# PostHouseImage Service

This is the PostHouseImage service

Generated with

```
micro new sss/PostHouseImage --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.PostHouseImage
- Type: srv
- Alias: PostHouseImage

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
./PostHouseImage-srv
```

Build a docker image
```
make docker
```