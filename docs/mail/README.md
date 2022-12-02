# ローカルでのメール操作

1. メールを送信する

```go:go
import "github.com/hack-31/point-app-backend/utils/email"

// 略

// メールを送信
result, err := email.SendMail("recipient@sample.com", "お知らせメール", "初めまして。\nお知らせです。\n")

if err != nil {
  // 失敗
  println(err.Error())
}
// 成功
println(result)
```

2. 送信したメールを確認

コンテナ内かホスト内かで、コマンドが異なる
また、ホスト側で確認する際は、[jq](https://formulae.brew.sh/formula/jq)を入れる必要がある

```sh
# ホスト
$ make read-mail-h
# コンテナ
$ make read-mail-c 
```

実行結果は以下となり、一番下がユーザに送られたメールの内容になる

```json
{
  "messages": [
    {
      "Id": "zjtkthbdnoztekij-yxymbmln-vbij-cvyi-mngp-qttznoskoozk-vwjnme",
      "Timestamp": "2022-12-02T17:19:07",
      "Region": "ap-northeast-1",
      "Source": "sender@sample.com",
      "Destination": {
        "ToAddresses": [
          "recipient@sample.com"
        ]
      },
      "Subject": "お知らせメール",
      "Body": {
        "text_part": "初めまして。\nお知らせです。\n",
        "html_part": null
      }
    }
  ]
}
```
