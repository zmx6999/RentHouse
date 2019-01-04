# UserAuth Service

This is the UserAuth service

Generated with

```
micro new sss/UserAuth --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.UserAuth
- Type: srv
- Alias: UserAuth

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
./UserAuth-srv
```

Build a docker image
```
make docker
```