# 🌍 OpenSource Globalizer AI

> Asistente de localización y mantenimiento impulsado por IA para proyectos de código abierto

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## ¿Qué es OpenSource Globalizer AI?

OpenSource Globalizer AI ayuda a los **mantenedores de código abierto a construir comunidades verdaderamente globales** mediante la automatización de:

- 📖 **Traducción de README / documentación** — conserva la estructura Markdown en más de 10 idiomas (mediante goldmark AST)
- 🔄 **Integración con GitHub Action** — traducción automática al hacer push, abre un PR automáticamente
- 🏷️ **Triaje y traducción de Issues** (V2) — detecta el idioma, etiqueta automáticamente, traduce para mantenedores no nativos
- 📦 **Generación de notas de versión** (V3) — produce notas de versión en varios idiomas a partir de changelogs

Todo impulsado por IA. Todo integrado en el flujo de trabajo de GitHub que ya usas.

---

## ¿Por qué?

El código abierto es global por naturaleza. Tus usuarios hablan 中文, 日本語, 한국어, Español, Français, Deutsch…
Pero la mayoría de los mantenedores no pueden traducir manualmente cada README, cada Issue, cada Nota de versión.

> **OpenSource Globalizer AI reduce la carga de localización de horas → segundos.**

---

## Características

| Característica | Ver | Estado | Descripción |
|---------|-----|--------|-------------|
| 📖 **Traductor de README** | v0.1 | ✅ Publicado | Traduce README.md a varios idiomas; goldmark AST conserva todo el formato |
| 🌐 **API HTTP** | v0.1 | ✅ Publicado | API REST mediante Gin, POST /api/v1/translate |
| 🔄 **Acción de GitHub** | v0.2 | ✅ Publicado | Traducción automática al hacer push, crear PR automáticamente |
| 🏷️ **Asistente de Issues** | v0.3 | ✅ Publicado | Detecta el idioma del issue, clasifica automáticamente, responde y etiqueta |
| 📦 **Asistente de versiones** | v0.4 | 📋 Planificado | Genera notas de versión en varios idiomas |
| 🤖 **App de GitHub** | v1.0 | 📋 Planificado | Integración completa del bot con comentarios y revisión de PR |

---

## Inicio rápido

> 📖 Guía de instalación completa: **[INSTALL.md](INSTALL.md)** — tutorial de 5 minutos.

### Instalación en una línea

```bash
# Download a prebuilt binary (macOS/Linux/Windows), or:
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker (no requiere entorno Go)

> **Requisito previo**: `-v $(pwd):/workspace` monta el directorio actual en el contenedor, así que **ejecútalo desde el directorio que contiene `README.md`**.
> macOS Docker Desktop solo comparte `/Users`, etc. de forma predeterminada; si ejecutas desde `/tmp` u otra ruta no compartida, el contenedor no verá tus archivos (`no such file or directory`).
> Solución: Docker Desktop → Settings → Resources → File Sharing → añade el directorio, o usa una ruta bajo `/Users/...`.

```bash
# Pull the image
docker pull ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0

# CLI mode: translate the README in the current directory (run from README.md's directory)
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0 \
  translate README.md --lang zh-CN,ja

# CLI mode: custom API endpoint and model
docker run --rm -e OPENAI_API_KEY="sk-xxx" \
  -v $(pwd):/workspace -w /workspace \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0 \
  translate README.md --lang zh-CN --base-url https://api.example.com/v1 -m gpt-5-mini

# HTTP API mode (no mount needed)
docker run -d -p 8080:8080 -e OPENAI_API_KEY="sk-xxx" \
  ghcr.io/ytc301/opensource-globalizer-ai:v0.3.0 serve
```


### Traducir en un comando

```bash
# Via environment variable
export OPENAI_API_KEY="sk-..."
globalizer translate README.md --lang zh-CN,ja,ko,es

# Via CLI flag
globalizer translate README.md --lang zh-CN,ja,ko,es --api-key "sk-..."

# Custom API endpoint and model
globalizer translate README.md --lang zh-CN -m gpt-5-mini \
  --base-url https://api.openai.com/v1 --api-key "sk-..."
```


```
📖 Source: README.md
🌍 Target languages: zh-CN, ja, ko, es
🤖 Model: gpt-4o
  ✅ docs/README.zh-CN.md
  ✅ docs/README.ja.md
  ✅ docs/README.ko.md
  ✅ docs/README.es.md
✨ Done! 4 language versions generated
```


### Iniciar API HTTP

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


### Acción de GitHub: traducción automática + PR automático (recomendado)

No se necesita entorno local: haz push de `README.md` y se crea automáticamente un PR de traducción:

```yaml
# .github/workflows/i18n.yml
name: AI Translation

on:
  push:
    branches: [main]       # tag pushes are excluded to keep PR creation reliable
    paths:
      - "README.md"
  workflow_dispatch:      # optional manual trigger

permissions:
  contents: write
  pull-requests: write

jobs:
  translate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: ytc301/OpenSource-Globalizer-AI/github-action@v0.2.0
        with:
          api-key: ${{ secrets.OPENAI_API_KEY }}
          languages: zh-CN,ja,ko
          model: gpt-4o                        # optional: custom model
          base-url: https://api.openai.com/v1  # optional: API-compatible endpoint
          output-dir: .                        # write README.<lang>.md at repo root
```


Los archivos traducidos se escriben en la raíz del repositorio como `README.<lang>.md` (p. ej., `README.zh-CN.md`, `README.ja.md`).
**Configuración:**

1. Añade el secreto `OPENAI_API_KEY`: repo **Settings → Secrets and variables → Actions**
2. Habilita la creación de PR: repo **Settings → Actions → General → Workflow permissions** → marca *Allow GitHub Actions to create and approve pull requests*
3. Haz push de `README.md` → la Action traduce automáticamente → se crea un PR (título `🌍 i18n: Auto-translate README to ...`)

> Soporta endpoints compatibles con OpenAI mediante las entradas `base-url` y `model`.
> ¿No tienes clave de API para probar? Deja `api-key` vacío y añade `mock: true` para verificar el flujo de principio a fin.
> El flujo de trabajo de ejemplo está en [.github/workflows/i18n.yml](.github/workflows/i18n.yml).

### Asistente de Issues: detección automática, clasificación, respuesta y etiquetado (v0.3)

El modo `serve` expone un endpoint de webhook de GitHub que procesa Issues automáticamente de principio a fin:

```bash
export GITHUB_TOKEN="ghp_xxx"                 # GitHub PAT with issues:write scope
export GLOBALIZER_WEBHOOK_SECRET="secret"     # webhook HMAC SHA-256 secret
globalizer serve                              # → POST /webhook is registered
```


**Configuración:**

1. Añade un webhook: repo **Settings → Webhooks → Add webhook**
   - **Payload URL**: `https://your-server/webhook`
   - **Content type**: `application/json`
   - **Secret**: el mismo valor que `GLOBALIZER_WEBHOOK_SECRET`
   - **Events**: **Issues** (`opened`, `edited`)
2. Configura `GITHUB_TOKEN` (repo **Settings → Secrets and variables → Actions**, o la variable de entorno del servidor) con el alcance `issues:write`

Cuando se abre un Issue que no está en inglés, el asistente automáticamente:

1. **Detecta** el idioma → añade una etiqueta `lang:xx`
2. **Lo clasifica** → añade una etiqueta `type:bug` / `type:feature` / `type:question` / `type:documentation`
3. **Publica** un resumen en inglés como primer comentario:

```
## 🌐 AI Translation

**语言:** zh-CN

**摘要:** Install fails on Ubuntu 24.04
```


> Las solicitudes de webhook se verifican con HMAC SHA-256 (`X-Hub-Signature-256`); las firmas no válidas se rechazan con `401`.
> La configuración se encuentra en `github.token` / `github.webhook_secret` en `.globalizer.yaml`, o en las variables de entorno `GITHUB_TOKEN` / `GLOBALIZER_WEBHOOK_SECRET`.
> El modelo de detección/clasificación usa `gpt-4o-mini` por defecto; configura `OPENAI_ISSUE_MODEL` para cambiarlo.

---

## Arquitectura

```
                 GitHub
                    |
            GitHub Action / CLI
                    |
        ┌───────────────────────┐
        |  OpenSource Globalizer |
        |────────────────────────|
        |  ┌─────────────────┐   |
        |  │  Gin HTTP API   │   |
        |  │  (serve cmd)    │   |
        |  ├─────────────────┤   |
        |  │  Translator     │───|── OpenAI API (GPT-4o)
        |  │  (goldmark AST) │   |
        |  ├─────────────────┤   |
        |  │  GORM + SQLite  │   |  ← translation cache
        |  └─────────────────┘   |
        └───────────────────────┘
```


Consulta [docs/architecture.md](docs/architecture.md) para ver el diseño completo.

---

## Pila tecnológica

| Capa | Tecnología |
|-------|-----------|
| **Lenguaje** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark (análisis a nivel de AST) |
| **ORM** | GORM + SQLite |
| **IA** | OpenAI API (GPT-4o / Codex) |
| **Configuración** | viper (fusión de env + YAML) |
| **Registro** | zap (estructurado) |
| **GitHub** | go-github, GitHub Actions |
| **Despliegue** | Docker, distribución de un solo binario |

---

## Estructura del proyecto

```
opensource-globalizer/
├── cmd/
│   └── globalizer/            # CLI entry (cobra + zap)
│       ├── main.go            # root command
│       ├── translate.go       # translate subcommand
│       ├── serve.go           # HTTP API server (Gin)
│       └── commands.go        # version / languages helpers
├── internal/
│   ├── handler/               # Gin HTTP handlers
│   ├── translator/            # translation engine (goldmark → AI → rebuild)
│   ├── ai/                    # AI Provider interface + OpenAI impl + Mock
│   ├── github/                # GitHub Client interface + Mock
│   └── store/                 # GORM + SQLite translation cache
├── pkg/
│   ├── markdown/              # goldmark AST parsing + segment management
│   └── config/                # viper configuration
├── docs/
│   ├── srs.md                 # Software Requirements Specification
│   ├── architecture.md        # architecture design
│   ├── api.md                 # API design
│   └── roadmap.md             # version roadmap
├── github-action/
│   └── action.yml             # GitHub Action definition
├── configs/
│   └── config.example.yaml    # config template
├── tests/
├── docker-compose.yml         # single-container deploy
├── Dockerfile                 # multi-stage build (Alpine)
├── Makefile
└── go.mod
```


---

## Hoja de ruta de versiones

| Versión | Cronograma | Entregable | Estado |
|---------|----------|-------------|--------|
| **v0.1.0** | 2026-07 (Semana 1-2) | Traductor de README CLI + API HTTP | ✅ Publicado |
| **v0.2.0** | 2026-07 (Semana 3-4) | Acción de GitHub + PR automático + Imagen Docker | ✅ Publicado |
| **v0.3.0** | 2026-08 | Detección de idioma de Issue + Clasificación + Respuesta automática + Etiqueta | ✅ Publicado |
| **v0.4.0** | 2026-09 | Generación de notas de versión en varios idiomas | 📋 Planificado |
| **v1.0.0** | 2026-10 | App de GitHub + Panel + Múltiples proveedores de IA | 📋 Planificado |

Consulta [docs/roadmap.md](docs/roadmap.md) para ver los hitos detallados.

---

## Contribuciones

¡Las contribuciones son bienvenidas! Por favor, lee primero [CONTRIBUTING.md](CONTRIBUTING.md).

```bash
git clone https://github.com/ytc301/OpenSource-Globalizer-AI.git
make deps
make test
make build
```


---

## Licencia

MIT © 2026 Colaboradores de OpenSource Globalizer AI