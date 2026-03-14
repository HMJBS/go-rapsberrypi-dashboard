# Issue草案（実装プランのフェーズ別）

このファイルは [implementation-plan.md](implementation-plan.md) の各フェーズを GitHub Issue として起票するための草案です。

現時点では、この草案をもとに GitHub Issue を起票済みです。
このファイルは、Issue 本文のテンプレート兼記録として保持します。

---

## P0: 設定JSON対応（フラグと併用）

**Title**: P0: Add JSON config file support (keep CLI flags override)

**Body**:

### Summary
運用時の引数地獄を避けるため、設定を JSON 1ファイルで渡せるようにする。既存のフラグ指定は後方互換として残し、JSONを上書きできるようにする。

### Goals
- `-config` フラグで JSON を読み込める
- JSON の値をベースに `cmd/dashboard/main.go` が `app.Config` を構築できる
- CLI フラグ指定があれば JSON を上書きできる（後方互換）
- 既定値・バリデーション（必須項目/範囲/単位）を定義する

### Non-goals
- 設定UIの追加
- 既存のフラグ削除（破壊的変更）

### Proposed design
- 新規: `internal/config` パッケージ
  - `Load(path string) (Config, error)`
  - `Validate()` / `ApplyDefaults()`
- 変更: `cmd/dashboard/main.go`
  - `-config` 追加
  - 読み込み順序: defaults → JSON → CLI flags override

### Acceptance criteria (DoD)
- `-config` 未指定でも従来通り動作する
- JSON指定時、指定した値が反映される
- 不正なJSON/値に対して分かりやすいエラーで落ちる
- 既存の `photos_dir` / `cache_dir` デフォルト挙動と矛盾しない

### Validation
- `mise run fmt` / `mise run lint` / `mise run vuln` / `mise run test`
- 可能なら `internal/config` にユニットテストを追加

### Files likely touched
- `cmd/dashboard/main.go`
- `internal/app/app.go`（必要なら構造体/型調整のみ）
- `internal/config/*`（新規）

---

## P0.5: デザイン案に合わせた仕様明文化（ドキュメント寄せ）

**Title**: P0.5: Align docs with designExample.html (source of truth)

**Body**:

### Summary
`design.png` はスクリーンショットであり、正本は [designExample.html](designExample.html)。このHTMLのレイアウト/表示要素に合わせて、仕様（features）を実装可能な粒度で固定する。

### Goals
- 仕様の正本を明確化（HTMLが正、PNGは参考）
- レイアウトを文章で固定（左:写真 / 右上:日時カード / 右下:天気カード）
- 表示要素（日時/天気/写真）を「必須/任意」に分ける
- 日付表示形式、天気の表示項目、文字列の英語表記ルールを定義する

### Non-goals
- UIの実装（このIssueはドキュメントのみ）
- 画像/フォントの追加（仕様確定が主）

### Proposed edits
- `docs/features.md`
  - 表示要素の再整理（HTMLに寄せる）
  - 時刻/日付/天気の表示フォーマットを明文化
- `docs/implementation-plan.md`
  - デザイン寄せフェーズの成果物（何が決まったら完了か）を明確に

### Acceptance criteria (DoD)
- features.md を読めば、UIの必須要素と表示フォーマットが判断できる
- `designExample.html` が正本であることが明記されている
- 実装者が「どこに何を表示するか」を迷わないレベルの記述がある

### Validation
- 文言の整合性チェック（features/architecture/run-on-raspberrypi との矛盾が無い）

### Files likely touched
- `docs/features.md`
- `docs/implementation-plan.md`

---

## P1: 日付表示（YYYY/MM/DD）

**Title**: P1: Render date (YYYY/MM/DD) in the datetime card

**Body**:

### Summary
デザイン（右上カード）に合わせて日付を表示する。時刻は既存のまま秒単位更新。

### Goals
- `YYYY/MM/DD` を表示できる
- タイムゾーン（`-tz`）が正しく反映される

### Non-goals
- フォントの大幅変更
- 多言語対応

### Proposed design
- `internal/app/app.go` の `render` で日付文字列を生成して `widgets.DrawText5x7` で描画
- タイムゾーンは `time.LoadLocation` + `now.In(loc)` に統一（必要なら）

### Acceptance criteria (DoD)
- 右上カード相当の領域に日付が表示される
- `-tz` を変えると日付/時刻が期待通り変わる

### Validation
- `-preview_dir` を使って `latest.png` で目視確認
- `mise run fmt` / `mise run lint` / `mise run vuln` / `mise run test`

### Files likely touched
- `internal/app/app.go`
- `internal/widgets/text5x7.go`（必要なら文字追加）

---

## P2: 天気表示の拡張（最高/最低・都市名・アイコン）

**Title**: P2: Improve weather card (icon, condition text, location, min/max)

**Body**:

### Summary
デザイン（右下カード）に合わせ、天気を「アイコン + 気温 + 状態文字列 + ロケーション」で表示できるようにする。APIは原則 Open-Meteo 継続（キー不要）。

### Goals
- Open-Meteo のレスポンス拡張で最高/最低（今日）を取得する
- 状態文字列（例: PARTLY CLOUDY）を表示できる
- ロケーション文字列（例: Tokyo, Japan）を表示できる
- アイコン（軽量描画）を表示できる

### Non-goals
- OpenWeatherMap への移行（必要性が出たら別Issue）
- 高解像度画像アイコンの導入

### Proposed design
- `internal/weather/weather.go`
  - `daily` から min/max を取得
  - 表示向けの `ConditionText()`（あるいはコード→英語表現マップ）を用意
- ロケーションはネットワーク依存を増やさないため、設定（P0 JSON/CLI）で指定する案を優先
- `internal/widgets` に簡易アイコン描画（雲/晴れ/雨/雪/雷など最小セット）
- `internal/app/app.go` の描画レイアウトをカード寄せ

### Acceptance criteria (DoD)
- 天気カードに「アイコン・気温・状態文・場所」が表示される
- 取得失敗時も、最後のキャッシュ表示/エラー表示で画面が壊れない
- 追加フィールドがあっても既存の `weather.json` が読み込める（or 自動更新される）

### Validation
- `-preview_dir` で目視
- `mise run fmt` / `mise run lint` / `mise run vuln` / `mise run test`

### Files likely touched
- `internal/weather/weather.go`
- `internal/app/app.go`
- `internal/widgets/*`（新規ファイル追加の可能性あり）

---

## P3: 画像まわりの堅牢性（任意）

**Title**: P3: Photo rendering robustness/performance improvements (optional)

**Body**:

### Summary
Pi 1B のリソース制約を踏まえ、画像ロード/リサイズがボトルネックになった場合に備えて改善する。

### Goals (pick as needed)
- 画面サイズに整形済みのキャッシュ（ディスク保存）を検討/実装
- 画像デコード失敗時の扱い（スキップ/次の画像へ）を明確化
- 画像が0枚の場合の表示を整備

### Non-goals
- Drive同期の内製化（rclone外出しのまま）

### Acceptance criteria (DoD)
- 明確な改善点（CPU/体感/安定性）が確認できる
- エラーが出ても落ちない/画面が壊れない

### Validation
- `-preview_dir` で長時間実行して挙動確認（可能なら）
- `mise run fmt` / `mise run lint` / `mise run vuln` / `mise run test`

### Files likely touched
- `internal/photos/photos.go`
- `internal/app/app.go`

---

## 起票手順（手動）

1. GitHub の `Issues` → `New issue` を開く
2. このファイルから該当フェーズの **Title** と **Body** をコピー
3. ラベル/マイルストーン（任意）を付けて作成

推奨起票順:
1. P0.5
2. P0
3. P1
4. P2
5. P3
