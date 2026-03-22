# Go テンプレート: 開発環境（Dev Environment）

このリポジトリは「Go 開発用テンプレート」として、ローカルと CI で同じ品質チェックを再現できる構成にしています。

コーディング規約は [docs/coding-standards.md](coding-standards.md) を参照してください。

## 目的

- **保存時**に自動フォーマット（VS Code）
- **コミット時**に自動フォーマット/静的解析（git hooks）
- **PR/Push**で必ずテスト・lint・脆弱性チェック（GitHub Actions CI）
- ツールや実行コマンドを `mise` に寄せて、環境差分を減らす

## 前提（ローカル）

- `mise` が入っていること
- Go 用エディタ（推奨: VS Code）

## 初期セットアップ

1) ツール導入（Go / golangci-lint2   など）

```bash
mise install
```

2) リポジトリ管理の git hooks を有効化

```bash
mise run hooks:install
```

これは `git config core.hooksPath .githooks` を設定し、コミット時のチェックを有効にします。

## 日常コマンド（mise / make）

- フォーマット

```bash
mise run fmt
# or
make fmt
```

- lint（静的解析）

```bash
mise run lint
# or
make lint
```

- テスト

```bash
mise run test
# or
make test
```

- ビルド（例: Linux armv6 向け）

```bash
mise run build
# or
make build
```

ビルド成果物は `bin/` 配下に出力されます。

## コミット時の自動チェック（pre-commit）

[/.githooks/pre-commit](../.githooks/pre-commit) で以下を実行します。

- ステージされた `.go` を `gofmt`（+ `goimports` があれば）で整形し、再ステージ
- `golangci-lint run ./...`
- `govulncheck ./...`

注意:
- `goimports` / `golangci-lint` / `govulncheck` は `mise install` で導入されます（[mise.toml](../mise.toml) の tools に指定）。
  - hooks 側は「PATH にある場合のみ実行」の実装にしてあるため、未導入でもコミット自体は可能です（ただし推奨は `mise install`）。
- 直接入れる場合は `mise use go:golang.org/x/vuln/cmd/govulncheck@latest` でもOKです。

## VS Code のライブ lint / 保存時フォーマット

[/.vscode/settings.json](../.vscode/settings.json) に設定があります。

- 保存時フォーマット: `goimports`（import 整理含む）
- 保存時 lint: `golangci-lint`（workspace モード）
- Language Server: `gopls`（`staticcheck` 有効）

推奨拡張は [/.vscode/extensions.json](../.vscode/extensions.json) に定義しています。

## CI（GitHub Actions）

[/.github/workflows/ci.yml](../.github/workflows/ci.yml) で PR/Push 時に次を実行します。

- `go test ./...`
- `golangci-lint`
- `govulncheck ./...`

ローカルで通っても CI で落ちる、またはその逆が起きないようにするのが狙いです。

## CHANGELOG 自動生成

Conventional Commits を前提に、Release Please で `CHANGELOG.md` を自動生成します。
詳細は [docs/git-workflow.md](git-workflow.md) を参照してください。

## Go バージョンの固定（再現性）

[go.mod](../go.mod) の `go` / `toolchain` を使って、使用する Go のメジャー/ツールチェインを揃えています。

## トラブルシュート

### FrameBuffer が無い開発機で動作確認したい

このアプリは本来 `/dev/fb0` に描画しますが、開発機に FrameBuffer が無い場合は
`-preview_dir` を指定すると、最新フレームを `latest.png` として書き出します（プレビューモード）。

例:

```bash
go run ./cmd/dashboard \
  -preview_dir ./out \
  -photos_dir /path/to/photos \
  -cache_dir  ./tmp/cache \
  -lat 35.681236 -lon 139.767125 -tz Asia/Tokyo
```

更新間隔を変えたい場合は `-preview_every` を指定します（例: `-preview_every 5s`）。

- `mise install` が TOML エラーで落ちる
  - `mise.toml` の tasks 名に `:` を含める場合、`[tasks."hooks:install"]` のように **クオートが必要**です。

- `golangci-lint` が設定を読めない（v2）
  - [/.golangci.yml](../.golangci.yml) に `version: 2` が必要です。
  - v2 では `gofmt/goimports` は `linters` ではなく `formatters` 側で有効化します。
