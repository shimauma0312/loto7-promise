# Loto7 Promise

### 結果取得
```bash
# 過去10回分
go run cmd/result/main.go 10

# 過去100回分
go run cmd/result/main.go 100
```

### ヒートマップ分析
```bash
# 分析 (過去50回)
go run cmd/heatmap/main.go

# 詳細分析
go run cmd/heatmap/main.go -range=100 -json

# 特定数字の分析
go run cmd/heatmap/main.go -number=7 -range=100
```