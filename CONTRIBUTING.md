# Contributing to OpenSource Globalizer AI

Thanks for your interest in contributing! 🎉

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/opensource-globalizer.git`
3. Create a branch: `git checkout -b feat/your-feature`
4. Make your changes
5. Run tests: `make test`
6. Commit using [Conventional Commits](https://www.conventionalcommits.org/): `feat: add xyz`
7. Push and open a Pull Request

## Development Setup

```bash
# Prerequisites
- Go 1.23+
- OpenAI API Key
- Docker (optional, only for serve mode)

# Run CLI mode (zero dependencies)
make run ARGS="translate README.md --lang zh-CN --dry-run"

# Run HTTP API mode
make run ARGS="serve"

# Run with Docker
make docker-build
make docker-up
```

## Commit Convention

We follow Conventional Commits:

- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation
- `refactor:` — Code refactoring
- `test:` — Tests
- `chore:` — Maintenance

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.23+ |
| CLI | cobra |
| HTTP | Gin |
| Markdown | goldmark |
| ORM | GORM + SQLite |
| AI | OpenAI API |
| Config | viper |
| Logging | zap |

## Code of Conduct

Be respectful. Be constructive. Help make open-source more global. 🌍

## Questions?

Open an Issue or start a Discussion.
