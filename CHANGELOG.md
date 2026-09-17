## Unreleased

## 1.0.0 (September 17, 2026)

FEATURES:

- Add `kilhog_network` and `kilhog_subnet` resources backed by the Kilhog Go SDK
- Add `kilhog_network` and `kilhog_subnet` data sources to read existing objects by ID or name
- Provider configuration via `base_url` and `api_key`

BREAKING CHANGES:

- `kilhog_network` tags are now a `map(string)` (`tags = { key = "value" }`) instead of a list of `{ key, value }` blocks. This fixes a permanent diff caused by the API returning tags in a different (sorted) order than they were written.
- Remove the `tags` attribute from the `kilhog_subnet` resource and data source. The Kilhog API does not persist tags on subnets, so the attribute always produced an "inconsistent result after apply" error.

ENHANCEMENTS:

- GitHub Actions CI: build, lint, unit tests, generate drift check, and acceptance tests against a spun-up Kilhog API
- Depend on the published `github.com/kilhog-io/kilhog` module instead of a machine-local `replace` path
- Publish each tagged release to the HashiCorp Terraform Registry (`kilhog-io/kilhog`) after the GitHub Release is created
