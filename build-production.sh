#!/bin/bash
# ===================================================================
# 本番環境用ビルドスクリプト
# ===================================================================
# 
# 使用方法:
#   chmod +x build-production.sh
#   ./build-production.sh
#
# 生成されるファイル:
#   - dist/loto7-api (APIサーバーバイナリ)
#   - dist/loto7-cli (CLI版バイナリ)
#   - dist/nginx.conf (nginx設定ファイル)
#   - dist/systemd/ (systemdサービスファイル)
#
# ===================================================================

set -e

echo "ロト7予想システム 本番用ビルド開始..."

# 変数設定
BUILD_DIR="dist"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GO_VERSION=$(go version | awk '{print $3}')

echo "ビルド情報:"
echo "   - バージョン: ${VERSION}"
echo "   - ビルド時刻: ${BUILD_TIME}"
echo "   - Go バージョン: ${GO_VERSION}"

# ビルドディレクトリの準備
echo "ビルドディレクトリを準備中..."
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}/systemd
mkdir -p ${BUILD_DIR}/scripts

# 依存関係の確認とダウンロード
echo "依存関係を確認中..."
go mod tidy
go mod verify

# リンターとテストの実行（オプション）
if command -v golangci-lint &> /dev/null; then
    echo "リンターを実行中..."
    golangci-lint run --timeout=5m || echo "リンターで警告がありました"
fi

echo "テストを実行中..."
go test ./... || {
    echo "テストが失敗しました。ビルドを中止します。"
    exit 1
}

# バイナリビルド（Linux用）
echo "Linux用バイナリをビルド中..."

# APIサーバー
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o ${BUILD_DIR}/loto7-api \
    ./cmd/api

# CLI版（既存のresultコマンド）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w" \
    -o ${BUILD_DIR}/loto7-cli \
    ./cmd/result

# バイナリに実行権限を付与
chmod +x ${BUILD_DIR}/loto7-api
chmod +x ${BUILD_DIR}/loto7-cli

# 設定ファイルのコピー
echo "設定ファイルをコピー中..."
cp nginx.conf ${BUILD_DIR}/

# systemdサービスファイルの生成
echo "systemdサービスファイルを生成中..."
cat > ${BUILD_DIR}/systemd/loto7-api.service << 'EOF'
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

# インストールスクリプトの生成
echo "インストールスクリプトを生成中..."
cat > ${BUILD_DIR}/scripts/install.sh << 'EOF'
#!/bin/bash
# ===================================================================
# ロト7予想API インストールスクリプト
# ===================================================================

set -e

# 設定値
INSTALL_DIR="/opt/loto7"
SERVICE_USER="loto7"
SERVICE_NAME="loto7-api"

echo "ロト7予想API インストール開始..."

# rootチェック
if [[ $EUID -ne 0 ]]; then
   echo "このスクリプトはroot権限で実行してください (sudo ./install.sh)"
   exit 1
fi

# ユーザー作成
echo "サービス用ユーザーを作成中..."
if ! id "${SERVICE_USER}" &>/dev/null; then
    useradd --system --home-dir ${INSTALL_DIR} --create-home --shell /bin/false ${SERVICE_USER}
    echo "ユーザー ${SERVICE_USER} を作成しました"
else
    echo "ユーザー ${SERVICE_USER} は既に存在します"
fi

# ディレクトリ作成
echo "インストールディレクトリを準備中..."
mkdir -p ${INSTALL_DIR}/logs
chown -R ${SERVICE_USER}:${SERVICE_USER} ${INSTALL_DIR}

# バイナリのコピー
echo "バイナリをコピー中..."
cp loto7-api ${INSTALL_DIR}/
cp loto7-cli ${INSTALL_DIR}/
chown ${SERVICE_USER}:${SERVICE_USER} ${INSTALL_DIR}/loto7-*
chmod +x ${INSTALL_DIR}/loto7-*

# systemdサービスの登録
echo "systemdサービスを登録中..."
cp systemd/loto7-api.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable ${SERVICE_NAME}

# nginx設定の提案
echo ""
echo "nginx設定ファイル (nginx.conf) がdistディレクトリにあります"
echo "   適切な場所にコピーして設定を調整してください"
echo ""

echo "インストール完了！"
echo ""
echo "サービス開始コマンド:"
echo "   sudo systemctl start ${SERVICE_NAME}"
echo ""
echo "サービス状態確認:"
echo "   sudo systemctl status ${SERVICE_NAME}"
echo ""
echo "ログ確認:"
echo "   sudo journalctl -u ${SERVICE_NAME} -f"
echo ""
EOF

# アンインストールスクリプトの生成
cat > ${BUILD_DIR}/scripts/uninstall.sh << 'EOF'
#!/bin/bash
# ===================================================================
# ロト7予想API アンインストールスクリプト
# ===================================================================

set -e

SERVICE_NAME="loto7-api"
SERVICE_USER="loto7"
INSTALL_DIR="/opt/loto7"

echo "ロト7予想API アンインストール開始..."

# rootチェック
if [[ $EUID -ne 0 ]]; then
   echo "このスクリプトはroot権限で実行してください (sudo ./uninstall.sh)"
   exit 1
fi

# サービス停止・無効化
echo "サービスを停止中..."
systemctl stop ${SERVICE_NAME} 2>/dev/null || true
systemctl disable ${SERVICE_NAME} 2>/dev/null || true

# systemdファイル削除
echo "systemdファイルを削除中..."
rm -f /etc/systemd/system/${SERVICE_NAME}.service
systemctl daemon-reload

# インストールディレクトリ削除
echo "インストールディレクトリを削除中..."
rm -rf ${INSTALL_DIR}

# ユーザー削除（確認付き）
read -p "サービス用ユーザー ${SERVICE_USER} を削除しますか？ (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    userdel ${SERVICE_USER} 2>/dev/null || true
    echo "ユーザー ${SERVICE_USER} を削除しました"
fi

echo "アンインストール完了！"
EOF

# スクリプトに実行権限を付与
chmod +x ${BUILD_DIR}/scripts/*.sh

# READMEの生成
echo "本番用READMEを生成中..."
cat > ${BUILD_DIR}/README_PRODUCTION.md << EOF
# ロト7予想API - 本番環境インストールガイド

## システム要件

- Linux (Ubuntu 20.04+ / CentOS 8+ / RHEL 8+ 推奨)
- systemd
- nginx (リバースプロキシ用)
- 最低 512MB RAM
- 最低 100MB ディスク容量

## インストール手順

### 1. ファイルの配置

```bash
# ビルド済みファイルを本番サーバーにコピー
scp -r dist/ user@production-server:/tmp/loto7-install/
```

### 2. インストール実行

```bash
# 本番サーバー上で
cd /tmp/loto7-install
sudo ./scripts/install.sh
```

### 3. サービス開始

```bash
# APIサーバー開始
sudo systemctl start loto7-api

# 自動起動有効化
sudo systemctl enable loto7-api

# 状態確認
sudo systemctl status loto7-api
```

### 4. nginx設定

```bash
# nginx設定をコピー
sudo cp nginx.conf /etc/nginx/sites-available/loto7-api
sudo ln -s /etc/nginx/sites-available/loto7-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 運用

### ログ確認

```bash
# リアルタイムログ
sudo journalctl -u loto7-api -f

# エラーログのみ
sudo journalctl -u loto7-api -p err
```

### サービス制御

```bash
# 開始
sudo systemctl start loto7-api

# 停止
sudo systemctl stop loto7-api

# 再起動
sudo systemctl restart loto7-api

# 状態確認
sudo systemctl status loto7-api
```

### 設定変更

環境変数を変更する場合は、systemdサービスファイルを編集：

```bash
sudo systemctl edit loto7-api
```

## アンインストール

```bash
cd /tmp/loto7-install
sudo ./scripts/uninstall.sh
```

## API動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# nginx経由でのアクセス
curl http://your-domain.com/api/v1/info

# 予想生成テスト
curl "http://your-domain.com/api/v1/prediction?use_trend=true"
```

## ⚠️ セキュリティ設定

1. ファイアウォール設定
2. nginx のセキュリティヘッダー設定
3. SSL/TLS証明書の設定
4. 定期的なアップデート

詳細は nginx.conf を参照してください。

---

ビルド情報:
- バージョン: ${VERSION}
- ビルド時刻: ${BUILD_TIME}
- Go バージョン: ${GO_VERSION}
EOF

# ビルド結果の表示
echo ""
echo "ビルド完了！"
echo ""
echo "生成されたファイル:"
echo "   - ${BUILD_DIR}/loto7-api (APIサーバーバイナリ)"
echo "   - ${BUILD_DIR}/loto7-cli (CLI版バイナリ)"
echo "   - ${BUILD_DIR}/nginx.conf (nginx設定)"
echo "   - ${BUILD_DIR}/systemd/loto7-api.service"
echo "   - ${BUILD_DIR}/scripts/install.sh"
echo "   - ${BUILD_DIR}/scripts/uninstall.sh"
echo "   - ${BUILD_DIR}/README_PRODUCTION.md"
echo ""
echo "バイナリサイズ:"
ls -lh ${BUILD_DIR}/loto7-* | awk '{print "   - " $9 ": " $5}'
echo ""
echo "次のステップ:"
echo "   1. distディレクトリを本番サーバーにコピー"
echo "   2. sudo ./scripts/install.sh を実行"
echo "   3. README_PRODUCTION.md の手順に従って設定"
echo ""
