# vaultpull

> A CLI tool to sync secrets from HashiCorp Vault into local `.env` files with profile support.

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Configure your Vault address and token, then pull secrets into a `.env` file:

```bash
# Set required environment variables
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your_token_here"

# Pull secrets from a Vault path into .env
vaultpull pull --path secret/data/myapp --out .env

# Use a named profile for different environments
vaultpull pull --profile staging --out .env.staging
```

**Example `.vaultpull.yaml` config:**

```yaml
profiles:
  production:
    path: secret/data/myapp/prod
    output: .env.production
  staging:
    path: secret/data/myapp/staging
    output: .env.staging
```

Then simply run:

```bash
vaultpull pull --profile production
```

---

## Commands

| Command | Description |
|---------|-------------|
| `pull` | Sync secrets from Vault to a local `.env` file |
| `profiles` | List available profiles from config |
| `version` | Print the current version |

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid Vault token with read access to the target paths

---

## License

[MIT](LICENSE) © 2024 yourusername