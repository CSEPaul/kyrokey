# kyrokey
`kyrokey` is a Go-based secret manager that writes secrets to the OS keychain and keeps a local SQLite index of `(service, username)` pairs so they can be listed and managed from one place.

## What the code does
The project provides two interfaces:
- **CLI mode** via `kyro k ...`
- **GUI mode** via `kyro g` (Fyne desktop app)

Core behavior:
- Stores secrets in the OS keyring (`github.com/zalando/go-keyring`)
- Tracks metadata (service + username only, never the secret) in `./conf/keychain_entries.db`
- Supports create, retrieve, list, delete-one-secret, and delete-tracking-db operations

## How it is structured
- `main.go` → app entrypoint
- `cmd/root.go` → root Cobra command (`kyro`)
- `cmd/kc_cli/*` → CLI commands (`set`, `get`, `list`, `del`, `deldb`)
- `cmd/kc_gui/*` → GUI screens wired to the same CLI business logic
- `libs/libs_keychain.go` → keyring + SQLite helper functions

## Build and run
From the project root:

- Run directly:
  - `go run . --help`
- Build binary:
  - `go build -o kyro .`
- Then use:
  - `./kyrokey --help`

## Command overview
Top-level commands:
- `kyrokey k` → CLI secret management commands
- `kyrokey g` → Launch GUI secret manager
- `kyrokey completion` → Generate shell completion scripts
- `kyrokey help` → Show command help

### CLI commands (`kyrokey k`)

#### 1) Set a secret
Purpose: Save a secret to keychain and register `(service, user)` in SQLite tracker.

Example:
- `kyrokey k set -s github -u paul -S my_token_value`

Flags:
- `-s, --service` (required): service name
- `-u, --user` (required): username/account
- `-S, --secret` (required): secret value

#### 2) Get a secret
Purpose: Read a secret from keychain for a given `(service, user)`.

Example:
- `kyrokey k get -s github -u paul`

Flags:
- `-s, --service` (required): service name
- `-u, --user` (required): username/account

#### 3) List tracked entries
Purpose: List all tracked `(service, user)` pairs from the local SQLite tracker.

Example:
- `kyrokey k list`

#### 4) Delete one secret
Purpose: Delete a secret from keychain and remove that `(service, user)` row from tracker DB.

Example:
- `kyrokey k del -s github -u paul`

Flags:
- `-s, --service`: service name
- `-u, --user`: username/account

#### 5) Delete tracking database
Purpose: Remove the local tracker DB file (`./conf/keychain_entries.db`) only.

Example:
- `kyrokey k deldb --Confirm Confirm`
- or: `kyrokey k deldb -C C`

Flags:
- `-C, --Confirm` (required): must be `Confirm` or `C` to proceed

## GUI mode
Run:
- `kyrokey g`

The GUI includes screens for:
- Set secret
- Get secret
- List tracked entries
- Delete a secret
- Delete tracking DB

## Typical workflow example
1. Save a secret:
   - `kyrokey k set -s aws -u deploy-user -S super-secret-value`
2. Verify it exists:
   - `kyrokey k list`
3. Retrieve it:
   - `kyrokey k get -s aws -u deploy-user`
4. Remove it when no longer needed:
   - `kyrokey k del -s aws -u deploy-user`

## Notes
- Secrets are stored in the OS keychain; this app only tracks service/user metadata in SQLite.
- `list` shows tracker entries for this app, not all keychain entries in the system.
