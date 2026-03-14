# Raspberry Pi 1B 上での実行（FrameBuffer + rclone）

このプロジェクトは X11 を使わず、Linux FrameBuffer（例: `/dev/fb0`）へ直接描画します。
Google Drive の画像はアプリが直接取得せず、`rclone` でローカルフォルダへ同期したものを参照します。

## 前提

- Raspberry Pi 1B
- FrameBuffer が利用できること（例: `/dev/fb0`）
- 画像同期に `rclone` を利用

## アプリの起動

例（写真は rclone 同期先、天気は緯度経度固定）:

```bash
./bin/dashboard-linux-armv6 \
	-config /etc/dashboard-config.json \
  -fb /dev/fb0 \
  -tz Asia/Tokyo
```

設定ファイル例:

```json
{
  "latitude": 35.681236,
  "longitude": 139.767125,
  "timezone": "Asia/Tokyo",
  "photos_dir": "/var/lib/dashboard/photos",
  "cache_dir": "/var/lib/dashboard/cache",
  "photo_interval_seconds": 60,
  "photo_rescan_interval_seconds": 300,
  "weather_interval_minutes": 10
}
```

既定では `/etc/dashboard-config.json` を読みに行きます。別パスを使う場合だけ `-config` を指定してください。フラグを併用した場合は、JSON の値よりフラグが優先されます。

## rclone で Drive → ローカル同期

### 初回セットアップ（例）

- 別マシンで `rclone config` を実行して Google アカウントの認証を済ませます
- 生成された `rclone.conf` を Pi へ配置します

Pi 側での例（パスは適宜変更）:

```bash
sudo mkdir -p /etc/rclone
sudo cp rclone.conf /etc/rclone/rclone.conf
sudo chmod 600 /etc/rclone/rclone.conf
sudo mkdir -p /var/lib/dashboard/photos
```

### 手動同期コマンド例

```bash
rclone --config /etc/rclone/rclone.conf sync \
  "gdrive:Pictures/dashboard" \
  /var/lib/dashboard/photos \
  --fast-list --delete-during
```

- `gdrive:` は rclone の remote 名です
- `Pictures/dashboard` は同期対象の Drive フォルダ（例）です

## 定期同期（Buildroot / SysVinit 想定）

Buildroot の軽量イメージでは systemd が入っていないことが多いため、
基本は BusyBox の `crond`（cron）か SysVinit の起動スクリプトで定期実行します。

### 方法A: cron（推奨）

1. `crond` を有効化して起動します（イメージ構成により手順は異なります）
2. root の crontab に同期コマンドを追加します

例（5分ごとに同期）:

```cron
*/5 * * * * /usr/bin/rclone --config /etc/rclone/rclone.conf sync "gdrive:Pictures/dashboard" /var/lib/dashboard/photos --fast-list --delete-during >/dev/null 2>&1
```

### 方法B: SysVinit の起動スクリプト + ループ

`/etc/init.d/S50dashboard-rclone-sync`（例）:

```sh
#!/bin/sh

case "$1" in
  start)
    (
      while true; do
        /usr/bin/rclone --config /etc/rclone/rclone.conf sync "gdrive:Pictures/dashboard" /var/lib/dashboard/photos --fast-list --delete-during
        sleep 300
      done
    ) &
    ;;
  stop)
    # 簡易実装: 必要なら pid 管理を追加する
    ;;
esac
```

※ この方式は簡単ですが、pid 管理やログなど運用面の作り込みが必要になりがちなので、可能なら cron を推奨します。

### 補足: systemd がある場合

systemd が利用できる環境では timer を使っても構いません。
