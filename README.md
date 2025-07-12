# gopkg

`gopkg` is a lightweight dependency manager for Go projects. It provides an alternative to Go modules by managing dependencies through a custom `gopkg.toml` file, with support for **both local and global installations**.

## Features

- Manage dependencies via `gopkg.toml`
- Supports **local** (`./gopkg_modules/`) and **global** (`~/.gopkg/modules/`) installation
- Lockfile support via `gopkg.lock`
- Auto detection mode (`--auto`) to scan Go imports and populate `gopkg.toml`
- Adds `replace` directives to `go.mod` automatically
- CLI commands: install, update, remove, check, list, versions
- Clean command to wipe installed modules, cache, and lockfile

## Installation

Install using `go install`:

```bash
go install github.com/pageton/gopkg@latest
```

Or build manually:

```bash
git clone https://github.com/pageton/gopkg
cd gopkg
go build -o gopkg .
```

## Usage

### 1. Initialize a project

```bash
gopkg init
```

Creates a `gopkg.toml` file and `gopkg_modules/` directory.

### 2. Add a dependency

```bash
gopkg add github.com/mattn/go-sqlite3@v1.14.17
```

To add it **globally**, use:

```bash
gopkg add -g github.com/mattn/go-sqlite3@v1.14.17
```

### 3. Install dependencies

```bash
gopkg install
```

Or using alias:

```bash
gopkg i
```

Install globally:

```bash
gopkg install --global
```

Use `--auto` to scan `.go` files and populate `gopkg.toml`:

```bash
gopkg install --auto
```

### 4. Update dependencies

```bash
gopkg update
```

Update a specific module:

```bash
gopkg update github.com/mattn/go-sqlite3@latest
```

Update globally:

```bash
gopkg update -g
```

### 5. Remove a dependency

```bash
gopkg remove github.com/mattn/go-sqlite3
```

Remove globally:

```bash
gopkg remove -g github.com/mattn/go-sqlite3
```

### 6. List available versions of a module

```bash
gopkg versions github.com/mattn/go-sqlite3
```

### 7. Check for outdated dependencies

```bash
gopkg check
```

### 8. Clean modules, lockfile, or cache

```bash
gopkg clean --lock --cache
```

## Project Structure

```
gopkg
├── cmd
│   ├── add.go
│   ├── check.go
│   ├── clean.go
│   ├── init.go
│   ├── install.go
│   ├── list.go
│   ├── remove.go
│   ├── root.go
│   ├── update.go
│   └── versions.go
├── core
│   ├── extract.go
│   ├── fetcher.go
│   ├── gomod.go
│   ├── importscan.go
│   ├── lockfile.go
│   ├── metadata.go
│   ├── module.go
│   ├── paths.go
│   └── version.go
├── go.mod
├── go.sum
└── main.go
```

## Contributing

Contributions are welcome. Please open an issue for suggestions or bugs, or submit a pull request directly. Make sure your code is clean and tested where appropriate.

## License

MIT License
