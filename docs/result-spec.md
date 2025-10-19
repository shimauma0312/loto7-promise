# ロト7結果取得

ロト7の過去の抽選結果を取得する機能。一度取得したデータはJSONで保存されて次回からは瞬時に読み込む。

## 処理フロー

### GetResult(N) の動作
```
1. キャッシュ状態チェック
   ├─ cache/metadata.json 読み込み
   └─ 最新回数・保存データ数確認

2. データ充足判定
   ├─ N ≤ 保存データ数 → Step6へ
   └─ N > 保存データ数 → Step3へ

3. 不足分データ取得
   ├─ calculateday.NewNumber() で最新回数取得
   ├─ 不足回数分のCSVをみずほ銀行から取得
   │   └─ 1秒間隔でHTTP GET
   └─ Shift-JIS → UTF-8 変換

4. データ解析・保存
   ├─ CSV行解析 (回数,数字1-7,ボーナス数字1-2)
   ├─ 既存データとマージ
   └─ JSON形式でcache/に保存

5. メタデータ更新
   ├─ 最終更新日時
   ├─ 最新回数
   └─ 総データ数

6. 結果返却
   └─ 最新N回分のデータを配列で返却
```

### ファイル構成
```
cache/
├─ loto7_results.json    # [{回数, 数字配列, 日付}, ...]
└─ metadata.json         # {last_update, latest_draw, total_count}
```

## API仕様

### 関数
```go
func GetResult(repeatNum int) [][]string
```

- **引数**: `repeatNum` - 取得したい過去回数（1以上）
- **戻り値**: 抽選結果の二次元配列（最新から古い順）

### 使用例
```bash
# 過去5回分
go run cmd/result/main.go 5

# 過去100回分（初回は時間がかかる）
go run cmd/result/main.go 100
```

### 出力形式
```
[[01 11 12 14 20 26 29] [03 08 15 22 25 31 37] ...]
```
