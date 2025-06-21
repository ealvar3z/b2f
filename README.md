# b2f (Bitwarden2Factotum)

A small Go utility to import your Bitwarden CSV-exported credentials into 9front's factotum.

## Features

* Parses Bitwarden CSV exports
* Extracts login URI, username, and password
* Generates `ctl` commands for factotum
* Optionally applies changes directly to `/mnt/factotum/ctl`

## Installation

1. Ensure you have [Go](https://golang.org/dl/) installed (Go 1.18+).
2. Clone this repository:

   ```sh
   git clone https://github.com/youruser/b2f.git
   cd b2f
   ```
3. Build the binary:

   ```sh
   go build -o b2f b2f.go
   ```

## Usage

```sh
# Read CSV and write commands to stdout:
./b2f -input bitwarden_export_YYYYMMDDHHMMSS.csv > facts.ctl

# Apply directly to factotum:
./b2f -input bitwarden_export_YYYYMMDDHHMMSS.csv -apply

# Use custom mountpoint or output file:
./b2f -input bitwarden_export_YYYYMMDDHHMMSS.csv -mount /path/to/factotum -out myfacts.ctl
```

### Flags

| Flag     | Default         | Description                                       |
| -------- | --------------- | ------------------------------------------------- |
| `-input` | *\<none>*       | Path to Bitwarden CSV export (required)           |
| `-mount` | `/mnt/factotum` | Factotum mountpoint                               |
| `-apply` | `false`         | Write directly to `<mount>/ctl` instead of stdout |
| `-out`   | *stdout*        | Path to write `ctl` commands                      |

## Examples

```sh
# Export from Bitwarden:
bw export --format csv --output bitwarden_export_YYYYMMDDHHMMSS.csv

# Generate factotum commands:
./b2f -input bitwarden_export_YYYYMMDDHHMMSS.csv > facts.ctl
cat facts.ctl > /mnt/factotum/ctl
```

## License

This project is licensed under the ISC License. See [LICENSE](LICENSE) for details.

