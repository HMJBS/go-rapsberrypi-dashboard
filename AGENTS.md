# Agent Instructions

このリポジトリを扱う自動化エージェント（AI/ボット含む）向けの指示です。ツール固有の前提は置きません。

## 原則

- 変更は目的に対して最小限にする（無関係な整形・大規模リファクタはしない）
- 既存の規約と自動化に従う（フォーマット/lint/脆弱性チェック）
- 追加の依存（Go module/外部ツール）は必要性を明確化し、最小にする
- 迷ったらドキュメントを優先する（仕様の追加・変更があるなら docs/README を更新）
- コメント・ドキュメントは日本語で書く（コード内のコメントも含む）

## 変更前に確認するドキュメント

- 開発環境: docs/dev-environment.md
- コーディング規約: docs/coding-standards.md
- レビュー観点: docs/review-checklist.md
- Git運用/リリース: docs/git-workflow.md

## ローカル検証（必須）

変更を入れたら、少なくとも以下を通すこと:

- `mise run fmt`
- `mise run lint`
- `mise run vuln`
- `mise run test`

## Git hooks / CI

- コミット時に `.githooks/pre-commit` が動く（未導入ツールがある場合はスキップする実装）
- PR title は CI で Conventional Commits を強制する（.github/workflows/pr-title.yml）
- `main` では Release Please が Release PR を作り、`CHANGELOG.md` を更新する（.github/workflows/release-please.yml）

## 出力・ログ・セキュリティ

- 秘密情報（トークン/鍵/個人情報）をログやドキュメントに書かない
- 既知の脆弱性や依存追加がある場合は、理由と影響範囲を明記する
