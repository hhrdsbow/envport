# envport

> CLI tool to snapshot and restore environment variable sets across projects

## Installation

```bash
go install github.com/yourusername/envport@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envport/releases).

## Usage

**Save a snapshot of your current environment:**
```bash
envport save myproject
```

**Restore a saved snapshot:**
```bash
envport load myproject
```

**List all saved snapshots:**
```bash
envport list
```

**Remove a snapshot:**
```bash
envport delete myproject
```

### Example Workflow

```bash
# Working on project A — save its environment
export DB_URL=postgres://localhost/projecta
export API_KEY=abc123
envport save project-a

# Switch to project B
envport load project-b

# Come back to project A later
envport load project-a
```

Snapshots are stored locally in `~/.envport/` as simple JSON files, making them easy to inspect or version control.

## License

[MIT](LICENSE)