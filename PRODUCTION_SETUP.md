# ロト7予想API - 本番環境セットアップ手順書

## 前提

- Linux サーバー (Ubuntu 20.04+)
- nginx
- systemd
- Go 1.23+

## 1. バイナリのビルド

### 開発環境でのビルド

```bash
# 依存関係の確認
go mod tidy

# テスト実行（オプション）
go test ./...

# 本番用バイナリをビルド
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o loto7-api ./cmd/api

# バイナリサイズ確認
ls -lh loto7-api
```

### 本番への転送

```bash
# バイナリとnginx設定を本番サーバーにコピー
scp loto7-api user@production-server:/tmp/
scp nginx.conf user@production-server:/tmp/
```

## 2. 本番セットアップ

### サービス用ユーザー作成

```bash
# loto7ユーザー作成
sudo useradd --system --home-dir /opt/loto7 --create-home --shell /bin/false loto7

# ディレクトリ作成
sudo mkdir -p /opt/loto7/logs
sudo chown -R loto7:loto7 /opt/loto7
```

### バイナリ配置

```bash
sudo cp /tmp/loto7-api /opt/loto7/
sudo chown loto7:loto7 /opt/loto7/loto7-api
sudo chmod +x /opt/loto7/loto7-api
```

### systemdサービスファイル作成

```bash
sudo tee /etc/systemd/system/loto7-api.service > /dev/null << 'EOF'
[Unit]
Description=Loto7 Prediction API Server
After=network.target
Wants=network.target

[Service]
Type=simple
User=loto7
Group=loto7
WorkingDirectory=/opt/loto7
ExecStart=/opt/loto7/loto7-api
Restart=always
RestartSec=5

# 環境変数
Environment=GIN_MODE=release
Environment=PORT=8080

# セキュリティ設定
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/loto7/logs
PrivateTmp=true

# リソース制限
LimitNOFILE=1024
MemoryMax=512M

# ログ設定
StandardOutput=journal
StandardError=journal
SyslogIdentifier=loto7-api

[Install]
WantedBy=multi-user.target
EOF
```

### サービスの有効化と開始

```bash
# systemd設定リロード
sudo systemctl daemon-reload

# サービス有効化
sudo systemctl enable loto7-api

# サービス開始
sudo systemctl start loto7-api

# ステータス確認
sudo systemctl status loto7-api
```

## 3. nginx設定

### nginx設定ファイルの配置

```bash
# 設定ファイルをコピー
sudo cp /tmp/nginx.conf /etc/nginx/sites-available/loto7-api

# シンボリックリンク作成
sudo ln -s /etc/nginx/sites-available/loto7-api /etc/nginx/sites-enabled/

# 設定ファイル編集（ドメイン名を本番用に変更）
sudo nano /etc/nginx/sites-available/loto7-api
```

### nginx設定の重要な変更点

```nginx
server {
    listen 80;
    # ↓ここを本番のドメイン名に変更
    server_name your-production-domain.com;
    
    # その他の設定はそのまま使用可能
}
```

### nginx設定テストと再起動

```bash
# 設定ファイルの構文チェック
sudo nginx -t

# nginx再起動
sudo systemctl reload nginx
```

## 4. 動作確認

### API直接アクセス

```bash
# ローカルでの動作確認
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/info
curl "http://localhost:8080/api/v1/prediction?use_trend=true"
```

### nginx経由でのアクセス

```bash
# nginx経由での動作確認
curl http://your-production-domain.com/health
curl http://your-production-domain.com/api/v1/info
curl "http://your-production-domain.com/api/v1/prediction?use_trend=true"
```

## 5. 運用コマンド

### サービス制御

```bash
# サービス状態確認
sudo systemctl status loto7-api

# ログ確認
sudo journalctl -u loto7-api -f

# サービス再起動
sudo systemctl restart loto7-api

# サービス停止
sudo systemctl stop loto7-api
```

### 設定変更

```bash
# 環境変数変更（ポート番号など）
sudo systemctl edit loto7-api

# 設定変更後の反映
sudo systemctl daemon-reload
sudo systemctl restart loto7-api
```

## 6. セキュリティ設定（推奨）

### ファイアウォール設定

```bash
# UFWの場合
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw enable

# 内部ポート8080は外部からアクセス不可にする
# （nginxが内部でプロキシするため）
```

### SSL/HTTPS設定（Let's Encrypt使用例）

```bash
# Certbot インストール
sudo apt install certbot python3-certbot-nginx

# SSL証明書取得
sudo certbot --nginx -d your-production-domain.com

# 自動更新設定確認
sudo certbot renew --dry-run
```

## 🗑️ 7. アンインストール手順

```bash
# サービス停止・無効化
sudo systemctl stop loto7-api
sudo systemctl disable loto7-api

# systemdファイル削除
sudo rm /etc/systemd/system/loto7-api.service
sudo systemctl daemon-reload

# nginx設定削除
sudo rm /etc/nginx/sites-enabled/loto7-api
sudo rm /etc/nginx/sites-available/loto7-api
sudo systemctl reload nginx

# インストールディレクトリ削除
sudo rm -rf /opt/loto7

# ユーザー削除（必要に応じて）
sudo userdel loto7
```

