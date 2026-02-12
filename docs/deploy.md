# Loto7 API デプロイ手順

本番環境へのデプロイ手順を説明します。

## 前提条件

- デプロイサーバーへのSSHアクセス権限
- Go 1.23以上がインストール済み
- systemdサービスが設定済み

## デプロイ手順

### 1. サーバーに接続

```bash
ssh user@server
```

### 2. プロジェクトディレクトリに移動

```bash
cd ~/develop/loto7-promise
```

### 3. 最新コードを取得

```bash
git checkout develop
git pull origin develop
```

### 4. ビルド

```bash
go build -o loto7-api ./cmd/api
```

### 5. サービスを停止

```bash
sudo systemctl stop loto7-api
```

### 6. バイナリとキャッシュを配置

```bash
sudo cp loto7-api /opt/loto7-api/loto7-api
sudo cp -r cache /opt/loto7-api/
```

### 7. サービスを起動

```bash
sudo systemctl start loto7-api
```

### 8. 動作確認

```bash
# サービスステータス確認
sudo systemctl status loto7-api

# ヘルスチェック
curl http://localhost:8080/health
```

正常に起動していれば、以下のようなレスポンスが返ります:

```json
{
  "success": true,
  "message": "サービスは正常に動作しています",
  "data": {
    "status": "healthy",
    "services": {
      "heatmap_engine": "operational",
      "random": "operational",
      "recommendation": "operational",
      "result_cache": "operational",
      "simulation": "operational"
    }
  }
}
```

## サービス管理コマンド

### ログ確認

```bash
# リアルタイムでログを監視
sudo journalctl -u loto7-api -f

# 最新のログ50行を表示
sudo journalctl -u loto7-api -n 50
```

### サービス再起動

```bash
sudo systemctl restart loto7-api
```

### サービスの自動起動確認

```bash
sudo systemctl is-enabled loto7-api
```

`enabled`と表示されれば、システム起動時に自動起動します。

## systemdサービス設定

サービスファイル: `/etc/systemd/system/loto7-api.service`

```ini
[Unit]
Description=Loto7 API Server
After=network.target

[Service]
Type=simple
User=tanari
Group=tanari
WorkingDirectory=/opt/loto7-api
ExecStart=/opt/loto7-api/loto7-api
Environment=GIN_MODE=release
Environment=PORT=8080
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

## トラブルシューティング

### サービスが起動しない時

1. ログを確認
   ```bash
   sudo journalctl -u loto7-api -n 100
   ```

2. バイナリの実行権限を確認
   ```bash
   ls -l /opt/loto7-api/loto7-api
   ```

3. 手動でバイナリを実行してエラーを確認
   ```bash
   cd /opt/loto7-api
   ./loto7-api
   ```

### ポートが既に使用されている時

```bash
# ポート8080を使用しているプロセスを確認
sudo lsof -i :8080
```

### サービス設定を変更した時

設定ファイルを変更した時は、systemdをリロードしてサービスを再起動:

```bash
sudo systemctl daemon-reload
sudo systemctl restart loto7-api
```
