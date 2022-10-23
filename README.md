# point-app-backend
ハッカソン2022で作成するポイントアプリのバックエンドプロジェクト

# 開発手順

```sh
# 環境手順
$ git clone https://github.com/hack-31/point-app-backend.git
$ cd ./point-app-backend
$ docker compose up -d --build
$ docker compose exec app sh
# サーバー起動
$ air
```

ホスト側で以下のURLでアクセス可能

API http://localhost:8081
swagger-ui http://localhost:80

ソースコードを書く際、vscodeを利用しているユーザは、[Dev container](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)を利用すると良い。
コンテナ内のファイルを直接修正でき、型補完や保存時フォーマットが使えるようになる。

今回、ホットリロードを利用しているため、保存のたびにサーバーを再起動する必要はない。
