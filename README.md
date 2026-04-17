# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/vaultpull/releases) page.

---

## Usage

```bash
# Authenticate and pull secrets from a Vault path into a local .env file
vaultpull --addr https://vault.example.com \
          --token $VAULT_TOKEN \
          --path secret/data/myapp \
          --out .env
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |
| `--path` | Secret path to pull from | *(required)* |
| `--out` | Output `.env` file path | `.env` |
| `--merge` | Merge with existing `.env` instead of overwriting | `false` |

**Example output (`.env`):**
```
DB_HOST=prod-db.internal
DB_PASSWORD=s3cr3t
API_KEY=abc123
```

> ⚠️ **Note:** `vaultpull` will never log or print secret values to stdout. Your secrets stay in the file.

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance (v1.10+)

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE) © 2024 yourusername