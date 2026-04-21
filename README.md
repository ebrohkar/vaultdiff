# vaultdiff

> CLI tool to diff and audit changes between HashiCorp Vault secret versions across environments

---

## Installation

```bash
go install github.com/yourusername/vaultdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultdiff.git
cd vaultdiff && go build -o vaultdiff .
```

---

## Usage

Compare two versions of a secret within the same Vault path:

```bash
vaultdiff --path secret/myapp/config --v1 3 --v2 5
```

Diff secrets across environments:

```bash
vaultdiff --src secret/staging/myapp --dst secret/production/myapp
```

Audit all changes to a secret over time:

```bash
vaultdiff --path secret/myapp/config --audit --since 2024-01-01
```

### Flags

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path |
| `--src` / `--dst` | Source and destination paths for cross-environment diff |
| `--v1` / `--v2` | Specific secret versions to compare |
| `--audit` | Show full version history with diffs |
| `--since` | Filter audit log from a given date |
| `--addr` | Vault server address (default: `$VAULT_ADDR`) |
| `--token` | Vault token (default: `$VAULT_TOKEN`) |

---

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine enabled

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)