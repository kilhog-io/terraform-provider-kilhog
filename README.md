# Terraform Provider Kilhog

Terraform provider for [Kilhog](https://kilhog.com), an IP Address Management (IPAM) platform. It uses the official Go SDK (`github.com/kilhog-io/kilhog/pkg/kilhog`) to manage networks and subnets.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.26
- A running Kilhog API (local or remote)

## Provider Configuration

```hcl
provider "kilhog" {
  base_url = "http://localhost:8080"
  api_key  = "your-api-key"
}
```

| Attribute | Description | Default |
|-----------|-------------|---------|
| `base_url` | Kilhog API base URL | `http://localhost:8080` or `KILHOG_BASE_URL` |
| `api_key` | Bearer token for authenticated APIs | `KILHOG_API_KEY` |

## Resources

| Resource | Description |
|----------|-------------|
| `kilhog_network` | Root tenancy container for subnets |
| `kilhog_subnet` | IPv4 subnet within a network or under a parent subnet |

## Data Sources

| Data Source | Description |
|-------------|-------------|
| `kilhog_network` | Read an existing network by ID or name |
| `kilhog_subnet` | Read an existing subnet by network ID and subnet ID or name |

### Example

```hcl
resource "kilhog_network" "production" {
  name        = "production"
  description = "Production network"
}

resource "kilhog_subnet" "dmz" {
  network_id  = kilhog_network.production.id
  name        = "dmz"
  address     = "10.0.0.0"
  prefix      = 24
}

resource "kilhog_subnet" "apps" {
  network_id       = kilhog_network.production.id
  parent_subnet_id = kilhog_subnet.dmz.id
  name             = "apps"
  prefix           = 25
}

data "kilhog_network" "production" {
  name = kilhog_network.production.name
}

data "kilhog_subnet" "dmz" {
  network_id = data.kilhog_network.production.id
  name       = "dmz"
}
```

See [`examples/`](examples/) for complete Terraform configurations.

## Local SDK Development

During local development, the provider depends on the Kilhog SDK via a `replace` directive in `go.mod`:

```go
replace github.com/kilhog-io/kilhog => /path/to/kilhog-cursor
```

## Building the Provider

```shell
go install
```

## Developing the Provider

```shell
make generate   # regenerate docs from schemas and examples
make test       # unit tests
make testacc    # acceptance tests (requires TF_ACC=1 and a running Kilhog API)
```

Acceptance tests use `KILHOG_BASE_URL` (default `http://localhost:8080`) and optional `KILHOG_API_KEY`.

```shell
TF_ACC=1 KILHOG_BASE_URL=http://localhost:8080 make testacc
```
