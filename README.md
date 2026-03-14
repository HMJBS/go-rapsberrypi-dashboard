# dashboard (Go)

個人的な趣味プロジェクトかつ今後のGo開発のテンプレートです。CI/CDやコード規約など、Goプロジェクトで一般的な構成を盛り込んでいます。

- 開発環境の詳細: [docs/dev-environment.md](docs/dev-environment.md)
- コーディング規約（チーム向け）: [docs/coding-standards.md](docs/coding-standards.md)
- レビュー観点（チェックリスト）: [docs/review-checklist.md](docs/review-checklist.md)
- Git 運用（チーム向け）: [docs/git-workflow.md](docs/git-workflow.md)
- CHANGELOG: [CHANGELOG.md](CHANGELOG.md)
- アーキテクチャ方針: [docs/architecture.md](docs/architecture.md)
- Raspberry Pi での実行（FrameBuffer + rclone）: [docs/run-on-raspberrypi.md](docs/run-on-raspberrypi.md)

## Prerequisites

- `mise` (tool + task runner)
- VS Code + Go extension (recommended)

## Setup

```bash
mise install
mise run hooks:install
```

## Common commands

```bash
mise run fmt
mise run lint
mise run test
mise run build
mise run vuln
```

Or via Make:

```bash
make fmt lint test build
```
