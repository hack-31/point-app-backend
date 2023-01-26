# リポジトリ層

- データストア(データベース)とドメインの橋かけを担う
- DB操作にまつわる操作(SQL)などは必ずリポジトリに記述する
- データの永続化と再構築を行う

# ファイル
## transaction.go
- トランザクションを扱うコード

```go:auth_router.go
// 利用方法
// router.go, auth_router.goでインスタンス化 
appConnection := repository.NewAppConnection(db)
// サービスに渡して利用
sendPointHandler := handler.NewSendPoint(&service.SendPoint{PointRepo: &rep, UserRepo: &rep, Connection: appConnection, DB: db})
```

### error.go
- エラーを定義する

## kvs.go
- Key/Values
- Redisにアクセスするための構造体

## repository.go
- DBのインスタンス作成やIFを定義
