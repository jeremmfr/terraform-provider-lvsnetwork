# terraform-provider-lvsnetwork

![GitHub release (latest by date)](https://img.shields.io/github/v/release/jeremmfr/terraform-provider-lvsnetwork)
[![Go Status](https://github.com/jeremmfr/terraform-provider-lvsnetwork/workflows/Go%20Tests/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvsnetwork/actions)
[![Lint Status](https://github.com/jeremmfr/terraform-provider-lvsnetwork/workflows/GolangCI-Lint/badge.svg)](https://github.com/jeremmfr/terraform-provider-lvsnetwork/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jeremmfr/terraform-provider-lvsnetwork)](https://goreportcard.com/report/github.com/jeremmfr/terraform-provider-lvsnetwork)

Terraform's provider for setup network interface and keepalived vrrp_instance on two server (master/slave) with [lvsnetwork-api](https://github.com/jeremmfr/lvsnetwork-api)

## Documentation

[terraform-provider-lvsnetwork](docs/index.md)  

Resources:

* [lvsnetwork_ifacevrrp](docs/resources/ifacevrrp.md)
* [lvsnetwork_vrrp_script](docs/resources/vrrp_script.md)

## Compile

```shell
go build
```
