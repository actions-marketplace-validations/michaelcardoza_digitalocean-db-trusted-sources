# DigitalOcean DB Trusted Sources

GitHub Action that adds or removes the runner's public IP from a DigitalOcean Managed Database's firewall ("trusted sources"). Lets a GitHub-hosted runner reach a locked-down managed database for the duration of a job.

## Why

You restricted your DigitalOcean Managed Database to trusted sources (good). Your CI/CD needs to run migrations / seed data / integration tests against that database. GitHub-hosted runners don't have stable IPs. This action whitelists the current runner's IP before your DB step and removes it after.

## Usage

```yaml
jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Whitelist runner IP
        uses: michaelcardoza/digitalocean-db-trusted-sources@v0
        with:
          version: 0.1.0
          action: add
          database_id: ${{ secrets.DO_DATABASE_ID }}
          digitalocean_token: ${{ secrets.DO_API_TOKEN }}

      - name: Run migrations
        run: ./scripts/migrate.sh

      - name: Remove runner IP
        if: always()                       # critical: don't leak IPs on failure
        uses: michaelcardoza/digitalocean-db-trusted-sources@v0
        with:
          version: 0.1.0
          action: remove
          database_id: ${{ secrets.DO_DATABASE_ID }}
          digitalocean_token: ${{ secrets.DO_API_TOKEN }}
```

**Always pair `add` (pre-job) with `remove` (post-job, `if: always()`).** Otherwise dead runner IPs accumulate in your firewall every CI run.

## Inputs

| Name | Required | Description |
|---|---|---|
| `version` | yes | Release version of the binary to download (e.g. `0.1.0`). Leading `v` optional. |
| `action` | yes | `add` or `remove`. |
| `database_id` | yes | UUID of the managed database. Get it from `doctl databases list`. |
| `digitalocean_token` | yes | DigitalOcean Personal Access Token with read+write database scope. |

Store both `database_id` and `digitalocean_token` as GitHub Secrets — never inline.

## Security

The action downloads the release binary and verifies its SHA-256 checksum before executing — a tampered binary fails to run. Releases are signed with Sigstore keyless (cosign) and ship SLSA build provenance attestations.

## Notes

- **Idempotent.** Running `add` when the IP is already trusted, or `remove` when it isn't, is a no-op.
- **Don't run concurrent jobs against the same database.** Two jobs adding/removing IPs at the same time can lose rules.
- **The `@vN` ref selects the action code; the `version` input selects the binary.** Pin the action to a tag (`@v0`) and set `version:` to the release whose binary you want. Keep them in step (action `@v0` + `version: 0.1.0`).

## License

[MIT](./LICENSE)
