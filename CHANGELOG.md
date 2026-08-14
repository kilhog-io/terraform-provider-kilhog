## 0.1.0 (Unreleased)

FEATURES:

- Add `kilhog_network` and `kilhog_subnet` resources backed by the Kilhog Go SDK
- Add `kilhog_network` and `kilhog_subnet` data sources to read existing objects by ID or name
- Provider configuration via `base_url` and `api_key`

ENHANCEMENTS:

- GitHub Actions CI: build, lint, unit tests, generate drift check, and acceptance tests against a spun-up Kilhog API
- Depend on the published `github.com/kilhog-io/kilhog` module instead of a machine-local `replace` path
- Publish each tagged release to the HashiCorp Terraform Registry (`kilhog-io/kilhog`) after the GitHub Release is created
