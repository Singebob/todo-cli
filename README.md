# todo-cli

A terminal UI for managing todos, organized by context.

Built with Go and [Charm](https://charm.sh) (Bubble Tea, Lip Gloss, Bubbles).

## Features

- Todos organized by **contexts** (work, perso, project...)
- Optional **categories** for filtering
- **Comments** on each todo (links, notes, context)
- Preview panel showing the selected todo's details
- Toggle done/undone, view completed todos per context
- Sorted oldest first
- Keyboard-driven, vim-style navigation

## Architecture

Hexagonal architecture — domain and application logic are decoupled from infrastructure:

```
internal/
  domain/         # Entities, repository ports, errors
  app/            # Service (use cases)
  infra/
    storage/      # JSON file adapter (~/.todo-cli/)
    inmemory/     # In-memory adapter (tests)
  ui/             # TUI adapter (Bubble Tea)
```

## Install

```bash
go install github.com/Singebob/todo-cli@latest
```

Or build from source:

```bash
git clone https://github.com/Singebob/todo-cli.git
cd todo-cli
go build -o todo-cli .
./todo-cli
```

## Keybindings

| Key | Action |
|---|---|
| `j` / `k` | Navigate todos |
| `Tab` / `Shift+Tab` | Switch context |
| `Enter` | Open todo detail |
| `Space` | Toggle done |
| `n` | New todo |
| `e` | Edit todo |
| `d` | Delete todo |
| `c` | Show/hide completed |
| `/` | Filter by category |
| `C` | Create context |
| `D` | Delete context |
| `q` | Quit |

## Storage

Todos and contexts are stored as JSON files in `~/.todo-cli/`.

## License

MIT
