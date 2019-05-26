# Inventory Service

This is the Inventory service

Generated with

```
micro new shop/inventory-srv --namespace=zw.com.shop --alias=inventory --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: zw.com.shop.srv.proto
- Type: srv
- Alias: proto

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend consul.

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
./inventory-srv
```

Build a docker image
```
make docker
```