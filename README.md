# dotze

> Keep your secrets quiet.

A zero-dependency CLI secret manager for developers. Encrypts your secrets locally with AES-256-GCM — no account, no server, no cloud required.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=flat)

---

## Why dotze?

- **`.env` files store secrets in plain text.** dotze encrypts everything with AES-256-GCM before touching disk.
- **No account, no server, no config.** Works fully offline. Your keys never leave your machine.
- **Share secrets safely.** Send an encrypted one-time link to a teammate — no Slack, no email, no plaintext.
- **Language agnostic.** Works with any runtime — Node, Python, Go, Ruby, anything.

---

## Installation

### go install

```bash
go install github.com/anto0102/dotze@latest
```

### From source

```bash
git clone https://github.com/anto0102/dotze
cd dotze
go build -ldflags="-s -w" -o dotze .
sudo mv dotze /usr/local/bin/
```

### Homebrew

```bash
brew install dotze  # coming soon
```

---

## Usage

### Initialize a vault

```bash
dotze init
```

Creates an encrypted `.dotze` vault in the current directory. The decryption key is stored in `~/.dotze/keys/` and never leaves your machine.

---

### Set a secret

```bash
dotze set API_KEY abc123
# ✓ API_KEY saved
```

---

### Get a secret

```bash
dotze get API_KEY
# abc123
```

Works with shell expansion:

```bash
export API_KEY=$(dotze get API_KEY)
```

---

### List all keys

```bash
dotze list
# API_KEY
# DATABASE_URL
# STRIPE_KEY

dotze list --show-values
# API_KEY=abc123
# DATABASE_URL=postgres://localhost/mydb
# STRIPE_KEY=sk_test_...
```

---

### Delete a secret

```bash
dotze delete API_KEY
# ✓ API_KEY deleted
```

---

### Run a command with secrets injected

```bash
dotze run -- node app.js
dotze run -- python main.py
dotze run -- go run main.go
```

All secrets are injected as environment variables into the subprocess. stdin/stdout/stderr are passed through transparently.

---

### Import from an existing `.env` file

```bash
dotze import           # reads .env in current directory
dotze import .env.prod # reads a specific file
# ✓ Imported 5 secrets from .env
```

---

### Share secrets with a teammate

dotze offers two ways to share secrets securely.

#### Via encrypted link (online)

```bash
dotze share
# ✓ Share link (expires in 24h):
#   https://0x0.st/abc.txt#password...
# ⚠  This link cannot be manually revoked.
```

The secrets are encrypted locally before upload. The decryption password lives in the URL fragment (`#password`) and is **never sent to the server**. The link is hosted on 0x0.st and expires automatically after 24 hours.

#### Via encrypted file (offline)

```bash
dotze share --local
# Enter password:
# Confirm password:
# ✓ Exported to secrets.enc

dotze share --local --out myteam.enc  # custom filename
```

Generates an encrypted file you can send via Slack, email, or AirDrop. The file is useless without the password.

---

### Pull secrets from a link or file

```bash
# From a share link
dotze pull https://0x0.st/abc.txt#password...
# Found 3 secrets: API_KEY, DATABASE_URL, STRIPE_KEY
# Import all? [y/N]: y
# ✓ Imported 3 secrets
# ℹ  Remote file will expire automatically.

# From an encrypted file
dotze pull secrets.enc
# Enter password:
# Found 3 secrets: API_KEY, DATABASE_URL, STRIPE_KEY
# Import all? [y/N]: y
# ✓ Imported 3 secrets from secrets.enc
```


## How it works

Each project gets a unique AES-256 key generated on `dotze init`. The key is stored in `~/.dotze/keys/<sha256-of-project-path>.key` and never leaves your machine. Secrets are encrypted using AES-256-GCM and stored in a `.dotze` file in your project directory.

When sharing, dotze generates a one-time AES-256 key, encrypts the secrets with it, uploads the ciphertext to [0x0.st](https://0x0.st), and embeds the decryption key in the URL fragment — a part of the URL that browsers and HTTP clients never send to servers by design. The remote file expires automatically after 24 hours.

---

## Roadmap

- [ ] `dotze git-hook install` — pre-commit hook that blocks plain `.env` files
- [ ] Homebrew distribution
- [ ] Team sync via S3 or git

---

## License

MIT