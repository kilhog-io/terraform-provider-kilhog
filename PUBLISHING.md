# Publishing to the HashiCorp Terraform Registry

This provider is released from GitHub tags. GoReleaser builds, signs, and
uploads the GitHub Release. After a **one-time** registration on
[registry.terraform.io](https://registry.terraform.io), HashiCorp installs a
webhook on this repository. Each new GitHub Release is then ingested as a
provider version automatically.

Published address:

```hcl
terraform {
  required_providers {
    kilhog = {
      source = "kilhog-io/kilhog"
    }
  }
}
```

Official reference: [Publish providers to the Terraform Registry](https://developer.hashicorp.com/terraform/registry/providers/publishing).

## Release flow

```
git tag vX.Y.Z  →  GitHub Actions (Release)
                 →  GoReleaser (signed zips + SHA256SUMS + manifest)
                 →  GitHub Release (published, not draft)
                 →  HashiCorp webhook (release event)
                 →  registry.terraform.io/providers/kilhog-io/kilhog
```

There is no separate HashiCorp API token for community providers. The registry
pulls artifacts from GitHub. Credentials needed in CI are only the GPG signing
key used to sign `SHA256SUMS`.

## One-time setup

### 1. Repository requirements

- Public GitHub repository named `terraform-provider-kilhog` (lowercase).
- Owned by the GitHub organization `kilhog-io` (this becomes the registry namespace).
- GitHub Actions allowed to run (org and repo settings). Third-party actions used
  by `.github/workflows/release.yml` must be permitted:
  - `crazy-max/ghaction-import-gpg`
  - `goreleaser/goreleaser-action`
- Docs generated under `docs/` (`make generate`) so the registry can render them.

### 2. Create an RSA GPG signing key

The registry **rejects ECC keys** (GnuPG default). Use RSA 4096.

```shell
gpg --full-generate-key
```

When prompted:

1. Kind of key: `1` (RSA and RSA)
2. Key size: `4096`
3. Expiry: `0` (does not expire) or an expiry you are willing to rotate
4. Real name / email: Kilhog release identity (for example `Kilhog Releases <security@kilhog.com>`)
5. Passphrase: strong, stored in a secrets manager (not in git)

List the key and note the fingerprint:

```shell
gpg --list-secret-keys --keyid-format=long
```

Export the **public** key (paste this into the Terraform Registry later):

```shell
gpg --armor --export "KEY_ID_OR_EMAIL"
```

Export the **private** key (GitHub Actions secret; keep it offline as well):

```shell
gpg --armor --export-secret-keys "KEY_ID_OR_EMAIL"
```

Back up the private key and passphrase. Losing them blocks new signed releases
until you register a replacement public key on the registry.

### 3. GitHub Actions secrets

In the repository: **Settings → Secrets and variables → Actions → New repository secret**.

| Secret | Value |
|--------|--------|
| `GPG_PRIVATE_KEY` | Full ASCII-armored private key, including `BEGIN` / `END` lines |
| `PASSPHRASE` | Passphrase for that private key |

`GITHUB_TOKEN` is provided by GitHub Actions. It does not need to be created.
The workflow uses `permissions: contents: write` to create the GitHub Release.

Do **not** put the GPG private key in Terraform Registry. Only the public key
goes there.

### 4. Publish the first GitHub Release

The registry will not accept a provider until at least one valid GitHub Release
exists. Tags must be Semantic Versioning with a `v` prefix. There must not be a
branch with the same name as the tag.

```shell
git checkout main
git pull origin main
git tag v0.1.0
git push origin v0.1.0
```

The **Release** workflow then:

1. Imports `GPG_PRIVATE_KEY`
2. Builds multi-platform binaries with GoReleaser
3. Attaches `terraform-provider-kilhog_<version>_manifest.json`
4. Attaches `terraform-provider-kilhog_<version>_SHA256SUMS` and `.sig`
5. Publishes a **final** (non-draft) GitHub Release

Confirm the release assets on GitHub before continuing. Never replace assets on
an already-published version; cut a new tag instead.

### 5. Sign in to the Terraform Registry

1. Open [registry.terraform.io](https://registry.terraform.io) and sign in with GitHub.
2. The GitHub account must be an **admin** of `kilhog-io/terraform-provider-kilhog`.
3. If the org uses OAuth App approval, an org owner must approve the
   **Terraform Registry** application.
4. Confirm scopes under GitHub **Settings → Applications → Authorized OAuth Apps → Terraform Registry**.

### 6. Register the GPG public key

1. In the Terraform Registry: **User Settings → Signing Keys** (or the
   organization signing keys if you publish under `kilhog-io`).
2. Add a new GPG key and paste the ASCII-armored **public** key from step 2.
3. Save. The key must belong to the same namespace you will publish under.

### 7. Declare the provider (once)

Publishing a provider to the public registry is **permanent**. You cannot
unpublish it. Double-check the GitHub org and repository.

1. In the Terraform Registry, open **Publish → Provider**.
2. Select organization `kilhog-io` and repository `terraform-provider-kilhog`.
3. Choose a category (for example **Cloud Tools** or **Network**).
4. Accept the [Registry Terms of Use](https://registry.terraform.io/terms).
5. Click **Publish**.

HashiCorp creates a GitHub webhook subscribed to `release` events. Future tags
that produce a GitHub Release are ingested without extra credentials.

### 8. Confirm the webhook

In GitHub: **Settings → Webhooks**. You should see a hook pointing at
`https://registry.terraform.io` (typically `/hooks/github`) for **Releases**.

If it is missing, open the provider page on the registry → **Settings → Resync**.
If several registry webhooks exist, delete the extras first, then Resync.

## Releasing later versions

After the one-time setup, each new version is:

```shell
# Update CHANGELOG.md, merge to main, then:
git tag vX.Y.Z
git push origin vX.Y.Z
```

The Release workflow publishes GitHub assets, then waits a few minutes for
`registry.terraform.io/v1/providers/kilhog-io/kilhog` to list the version.

If the provider is not registered yet, that wait step succeeds with a notice so
the first GitHub Release can still be created (required before step 7).

## Troubleshooting

| Symptom | What to check |
|---------|----------------|
| Release workflow fails at **Import GPG key** | `GPG_PRIVATE_KEY` / `PASSPHRASE` secrets; RSA key, not ECC |
| Workflow fails at GoReleaser signing | Fingerprint / passphrase; `--batch` signing is configured in `.goreleaser.yml` |
| GitHub Release missing `_SHA256SUMS.sig` or `_manifest.json` | Workflow asset check failed; inspect GoReleaser logs |
| Registry page missing the new version | Webhook Recent Deliveries; **Resync** on the provider settings page |
| `terraform init` signature error | Public key on the registry must match the private key used in CI |
| Tag did not trigger the workflow | Tag must match `vMAJOR.MINOR.PATCH` (optional prerelease suffix) |

Support: [terraform-registry@hashicorp.com](mailto:terraform-registry@hashicorp.com).
