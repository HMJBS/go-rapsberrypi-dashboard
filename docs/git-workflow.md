# Git 運用（チームテンプレート）

このリポジトリは **GitHub Flow** を前提に運用します。
GitHub の UI 上では "PR (Pull Request)" が正式名称ですが、MR（Merge Request）という呼び方でも同義として扱います。

## 運用ルール（必須）

### 1) GitHub Flow を使う

- 常に `main` をデプロイ可能（安定）な状態に保つ
- 作業は **短命ブランチ**（feature/fix など）で行い、PR を通して `main` に戻す
- `main` への直接 push はしない（保護ブランチで防ぐ）

基本手順:

1. `main` からブランチ作成
2. 変更をコミットして push
3. PR 作成（説明と動作確認を記述）
4. Reviewer を **必ず1名以上**アサイン
5. CI が緑・承認済みになったら **Squash Merge**
6. マージ後はブランチ削除

### 2) PR/MR には必ず Reviewer を 1 人アサインする

- PR を作った人は、**必ず** Reviewer を 1 名以上指定する
- Reviewer は内容に応じて適切な人を選ぶ（迷ったらチームの当番/ローテに従う）

推奨: GitHub の Branch protection で「Require approvals = 1」を必須化する（後述）。

### 3) PR/MR のメッセージは Conventional Commits に準拠する

このテンプレでは **Squash Merge** を使うため、最終的に `main` に入るコミットメッセージ（= PR title / squash メッセージ）が重要です。

- PR title（推奨）: Conventional Commits 形式
- Squash merge の commit message: PR title をベースに Conventional Commits になるように整える

フォーマット:

```
<type>(<scope>)?: <subject>

<body>

<footer>
```

例:

- `feat: add dashboard build target`
- `fix(ci): run govulncheck in workflow`
- `docs: add review checklist`
- `refactor!: split main into cmd/dashboard`（破壊的変更がある場合は `!` を付ける）

主な type（例）:
- `feat` / `fix` / `docs` / `refactor` / `test` / `chore` / `ci` / `build`

注意:
- subject は命令形・短く（末尾ピリオド不要）
- scope は任意（複数パッケージにまたがる場合は省略でもよい）

#### PR title の自動チェック（必須）

このテンプレでは PR title が Conventional Commits 形式になっているかを CI で検証します。

- ワークフロー: [/.github/workflows/pr-title.yml](../.github/workflows/pr-title.yml)
- draft PR は対象外（Ready for review になるとチェックされます）

推奨: GitHub の Branch protection で、このチェックを必須にしてください。
（Settings → Branches → Branch protection rules → Require status checks to pass → `pr-title`）

### 4) Squash Merge を使う

- PR は **Squash and merge** で `main` に取り込む
- `main` の履歴を読みやすく保ち、"1 PR = 1 変更" を原則にする

運用:
- PR のコミットが複数でも OK（作業の途中経過を残してよい）
- ただし squash 後の最終メッセージは Conventional Commits に整える

## ブランチ命名（推奨）

- `feat/<short-description>`
- `fix/<short-description>`
- `chore/<short-description>`
- `docs/<short-description>`

例:
- `feat/add-armv6-build`
- `fix/ci-govulncheck`

## PR の書き方（推奨）

- PR テンプレート: [/.github/pull_request_template.md](../.github/pull_request_template.md)
- 重要: 何を/なぜ/どう確認したか を短く書く
- 可能なら PR は小さく（レビューしやすさ優先）

## GitHub 設定（リポジトリ設定での推奨）

テンプレとして使う場合、GitHub 側の設定も揃えると運用がブレません。

### Branch protection（`main`）

推奨設定:
- Require a pull request before merging
- Require approvals: **1**（チーム規模に応じて 2 でも可）
- Require status checks to pass: **on**（例: `ci`）
- Require conversation resolution: **on**
- Do not allow bypassing the above settings

### Merge methods

- Allow squash merging: **on**
- Allow merge commits: **off**（推奨）
- Allow rebase merging: 好みだが、Squash運用なら **off** でもOK

## よくある判断基準

- 「小さな改善だけ」でも PR を切る（`main` 直pushを避ける）
- 急ぎでも Reviewer は必須（最短レビューを依頼する）
- Conventional Commits に迷ったら `chore:` を避け、できるだけ意図に近い type を選ぶ

## CHANGELOG / リリース（自動化）

このテンプレでは Conventional Commits を前提に、Release Please で `CHANGELOG.md` とリリースを自動化します。

- 設定: [/release-please-config.json](../release-please-config.json) / [/.release-please-manifest.json](../.release-please-manifest.json)
- ワークフロー: [/.github/workflows/release-please.yml](../.github/workflows/release-please.yml)

### 仕組み（Release PR 方式）

- `main` に Conventional Commits が入るたびに、Release Please が **Release PR** を作成/更新します
- Release PR には、次の内容が自動で含まれます
	- `CHANGELOG.md` の更新
	- バージョン更新（manifest）
- Release PR をマージすると、タグ（例: `v1.2.3`）と GitHub Release が作成されます

### CHANGELOG のポリシー

- `feat`/`fix`/`perf` を中心に記載
- `docs`/`ci`/`chore` などは原則として CHANGELOG に載せません（設定で hidden）

### 注意（Squash merge と Conventional Commits）

Release Please は `main` のコミットメッセージを解析します。
Squash merge 時の最終コミットメッセージが Conventional Commits になるように、次を徹底します。

- PR title を Conventional Commits にする
- Squash merge のデフォルトメッセージは **PR title ベース**にする（運用で統一）
