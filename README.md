# envy

Securely store and sync `.env` files across machines using Cloudflare R2. Files are encrypted locally before upload — your secrets never leave your machine unencrypted.

## How it works

1. On first run, `envy init` generates a local encryption key stored at `~/.envy/config.json`
2. `envy push` encrypts your `.env` file and uploads it to R2
3. `envy sync` downloads and decrypts it into any directory on any machine

The encryption key never leaves your machine. Without it, the files in R2 are unreadable.

## Install

Requires Go 1.21+.

```bash
go install github.com/pradeepbgs/envy@latest
```

Make sure `~/go/bin` is in your PATH:

```bash
# add to ~/.zshrc or ~/.bashrc
export PATH="$PATH:$HOME/go/bin"
```

## Setup

You'll need a [Cloudflare R2](https://developers.cloudflare.com/r2/) bucket and an API token with read/write permissions.

```bash
envy init
```

You'll be prompted for:
- R2 Endpoint — `https://<account_id>.r2.cloudflarestorage.com`
- Access Key ID
- Secret Access Key
- Bucket name

Config is saved to `~/.envy/config.json`. Keep this file safe — losing it means losing access to your encrypted envs.

## Usage

### Push a `.env` file to R2

```bash
envy push <name> <path-to-env-file>
```

```bash
envy push myapp /path/to/project/.env
```

### Sync (download) a `.env` file

```bash
envy sync <name> <target-directory>
```

```bash
envy sync myapp /path/to/project
```

This writes the decrypted `.env` to `<target-directory>/.env`. Use `--force` to overwrite an existing file:

```bash
envy sync myapp /path/to/project --force
```

### List all stored envs

```bash
envy list
```

### Delete an env from R2

```bash
envy delete <name>
```

## Example workflow

```bash
# machine 1 — push your env
envy push api-service /home/user/projects/api/.env

# machine 2 — pull it down
envy sync api-service /home/user/projects/api
```

---

## Upcoming features

- **Smarter sync path** — currently `sync` requires the full directory path. Next version will let you pass just a folder or file name relative to your current terminal directory, so `envy sync myapp .` will just work.

## License

MIT
