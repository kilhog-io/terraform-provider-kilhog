# Terraform Provider Kilhog

Terraform provider for [Kilhog](https://kilhog.com), built with the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building the Provider

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).

```shell
go get github.com/author/dependency
go mod tidy
```

## Using the Provider

```hcl
provider "kilhog" {
  endpoint = "https://api.kilhog.com"
}
```

## Developing the Provider

To compile the provider, run `go install`. This builds the provider and puts the binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
