# 🌍 OpenSource Globalizer AI

> Assistant de localisation et de maintenance par IA pour projets open-source

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP-informational)](#-roadmap)

---

## Qu'est-ce qu'OpenSource Globalizer AI ?

OpenSource Globalizer AI aide **les mainteneurs open-source à bâtir des communautés véritablement mondiales** en automatisant :

- 📖 **Traduction des README / documentation** — préserve la structure Markdown dans plus de 10 langues (via l'AST de goldmark)
- 🔄 **Intégration GitHub Action** — traduction automatique dès le push, ouverture automatique d'une PR
- 🏷️ **Tri et traduction des issues** (V2) — détecte la langue, étiquette automatiquement, traduit pour les mainteneurs non natifs
- 📦 **Génération des notes de version** (V3) — produit des notes de version multilingues à partir des changelogs

Le tout propulsé par l'IA. Le tout intégré au workflow GitHub que vous utilisez déjà.

---

## Pourquoi ?

L'open-source est mondial par nature. Vos utilisateurs parlent 中文, 日本語, 한국어, Español, Français, Deutsch…
Mais la plupart des mainteneurs ne peuvent pas traduire manuellement chaque README, chaque issue, chaque note de version.

> **OpenSource Globalizer AI réduit la charge de travail de localisation de quelques heures à quelques secondes.**

---

## Fonctionnalités

| Fonctionnalité | Ver | Statut | Description |
|----------------|-----|--------|-------------|
| 📖 **Traducteur de README** | v0.1 | ✅ Publié | Traduit le README.md en plusieurs langues, l'AST de goldmark préserve toute la mise en forme |
| 🌐 **API HTTP** | v0.1 | ✅ Publié | API REST via Gin, POST /api/v1/translate |
| 🔄 **Action GitHub** | v0.2 | ✅ Publié | Traduction automatique lors du push, création automatique de PR |
| 🏷️ **Assistant d'issues** | v0.3 | ✅ Publié | Détecte la langue de l'issue, classification automatique, réponse automatique + étiquette |
| 📦 **Assistant de version** | v0.4 | 📋 Planifié | Génère des notes de version multilingues |
| 🤖 **Application GitHub** | v1.0 | 📋 Planifié | Intégration complète du bot avec commentaires et revue de PR |

---

## Démarrage rapide

> 📖 Guide d'installation complet : **[INSTALL.md](INSTALL.md)** — présentation en 5 minutes.

### Installation en une ligne

```bash
# Download a prebuilt binary (macOS/Linux/Windows), or:
go install github.com/ytc301/opensource-globalizer/cmd/globalizer@latest
```


### Docker (aucun environnement Go requis)

> **Prérequis** : `-v $(pwd):/workspace` monte le répertoire courant dans le conteneur, alors **exécutez-le depuis le répertoire contenant `README.md`**.
> Docker Desktop pour macOS ne partage que `/Users`, etc. par défaut — si vous exécutez depuis `/tmp` ou un autre chemin non partagé, le conteneur ne verra pas vos fichiers (`no such file or directory`).
> Correctif : Docker Desktop → Paramètres → Ressources → Partage de fichiers → ajoutez le répertoire, ou utilisez un chemin sous `/Users/...`.

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


### Traduire en une commande

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


### Démarrer l'API HTTP

```bash
globalizer serve
# → curl -X POST http://localhost:8080/api/v1/translate \
#     -H 'Content-Type: application/json' \
#     -d '{"content":"# Hello","target_langs":["zh-CN"]}'
```


### Action GitHub : traduction automatique + PR automatique (recommandé)

Aucun environnement local requis — poussez `README.md` et une PR de traduction est créée automatiquement :

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


Les fichiers traduits sont écrits à la racine du dépôt sous la forme `README.<lang>.md` (par ex. `README.zh-CN.md`, `README.ja.md`).
**Configuration :**

1. Ajoutez le secret `OPENAI_API_KEY` : dépôt **Paramètres → Secrets et variables → Actions**
2. Activez la création de PR : dépôt **Paramètres → Actions → Général → Permissions des workflows** → cochez *Autoriser GitHub Actions à créer et approuver des pull requests*
3. Poussez `README.md` → l'Action traduit automatiquement → une PR est créée (titre `🌍 i18n: Auto-translate README to ...`)

> Prend en charge les points de terminaison compatibles OpenAI via les entrées `base-url` et `model`.
> Pas de clé API pour tester ? Laissez `api-key` vide et ajoutez `mock: true` pour vérifier le flux de bout en bout.
> Un exemple de workflow se trouve dans [.github/workflows/i18n.yml](.github/workflows/i18n.yml).

### Assistant d'issues : détection, classification, réponse et étiquetage automatiques (v0.3)

Le mode `serve` expose un point de terminaison webhook GitHub qui traite automatiquement les issues de bout en bout :

```bash
export GITHUB_TOKEN="ghp_xxx"                 # GitHub PAT with issues:write scope
export GLOBALIZER_WEBHOOK_SECRET="secret"     # webhook HMAC SHA-256 secret
globalizer serve                              # → POST /webhook is registered
```


**Configuration :**

1. Ajoutez un webhook : dépôt **Paramètres → Webhooks → Ajouter un webhook**
   - **URL de payload** : `https://your-server/webhook`
   - **Type de contenu** : `application/json`
   - **Secret** : même valeur que `GLOBALIZER_WEBHOOK_SECRET`
   - **Événements** : **Issues** (`opened`, `edited`)
2. Définissez `GITHUB_TOKEN` (dépôt **Paramètres → Secrets et variables → Actions**, ou l'environnement du serveur) avec le champ d'action `issues:write`

Lorsqu'une issue non anglophone est ouverte, l'assistant automatiquement :

1. **Détecte** la langue → ajoute une étiquette `lang:xx`
2. **La classe** → ajoute une étiquette `type:bug` / `type:feature` / `type:question` / `type:documentation`
3. **Publie** un résumé en anglais comme premier commentaire :

```
## 🌐 AI Translation

**语言:** zh-CN

**摘要:** Install fails on Ubuntu 24.04
```


> Les requêtes webhook sont vérifiées avec HMAC SHA-256 (`X-Hub-Signature-256`) ; les signatures invalides sont rejetées avec `401`.
> La configuration se trouve sous `github.token` / `github.webhook_secret` dans `.globalizer.yaml`, ou via les variables d'environnement `GITHUB_TOKEN` / `GLOBALIZER_WEBHOOK_SECRET`.
> Le modèle de détection/classification est `gpt-4o-mini` par défaut ; définissez `OPENAI_ISSUE_MODEL` pour le remplacer.

---

## Architecture

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


Voir [docs/architecture.md](docs/architecture.md) pour la conception complète.

---

## Pile technique

| Couche | Technologie |
|--------|-------------|
| **Langage** | Go 1.23+ |
| **CLI** | cobra |
| **HTTP** | Gin |
| **Markdown** | goldmark (analyse au niveau de l'AST) |
| **ORM** | GORM + SQLite |
| **IA** | API OpenAI (GPT-4o / Codex) |
| **Configuration** | viper (fusion env + YAML) |
| **Journalisation** | zap (structurée) |
| **GitHub** | go-github, GitHub Actions |
| **Déploiement** | Docker, distribution en binaire unique |

---

## Structure du projet

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

## Feuille de route des versions

| Version | Calendrier | Livrable | Statut |
|---------|------------|----------|--------|
| **v0.1.0** | 2026-07 (Semaine 1-2) | Traducteur README CLI + API HTTP | ✅ Publié |
| **v0.2.0** | 2026-07 (Semaine 3-4) | Action GitHub + PR automatique + Image Docker | ✅ Publié |
| **v0.3.0** | 2026-08 | Détection de langue des issues + Classification + Réponse automatique + Étiquette | ✅ Publié |
| **v0.4.0** | 2026-09 | Génération multilingue des notes de version | 📋 Planifié |
| **v1.0.0** | 2026-10 | Application GitHub + Tableau de bord + Multi-fournisseur d'IA | 📋 Planifié |

Voir [docs/roadmap.md](docs/roadmap.md) pour les jalons détaillés.

---

## Contribuer

Les contributions sont les bienvenues ! Merci de lire d'abord [CONTRIBUTING.md](CONTRIBUTING.md).

```bash
git clone https://github.com/ytc301/OpenSource-Globalizer-AI.git
make deps
make test
make build
```


---

## Licence

MIT © 2026 Contributeurs OpenSource Globalizer AI