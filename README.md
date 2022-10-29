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

[Makefile](./Makefile)Makefileを擁しているため以下のコマンドでも可能

```sh
$ make build-up
$ make serve
```

詳しくは、makeコマンドを実行

サーバー起動したら、ホスト側で以下のURLでアクセス可能

- API http://localhost:8081
- swagger-ui http://localhost:80
- adminer http://localhost:8082

ソースコードを書く際、vscodeを利用しているユーザは、[Dev container](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)を利用して開発することを推奨する。
コンテナ内のファイルを直接修正でき、[Goの拡張機能](https://github.com/golang/vscode-go)により型補完や保存時フォーマットが使えるようになる。

dev containerでの開発方法
1. vscodeに[Dev container](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)を導入
2. cmd+shift+pを押し、`>open folder in container`を押し、このプロジェクトを選択
3. 右下の通知に表示されるGo拡張機能に必要なモジュールをインストールボタンを押し、インストールする
4. vscode内のターミナルで`$ air`でサーバー起動

今回、ホットリロードを利用しているため、保存のたびにサーバーを再起動する必要はない。

# DB
開発用データベース情報

|項目|値|
|---|---|
|データベース種類|MySQL|
|サーバ|db|
|ユーザ名|admin|
|パスワード|password|
|データベース|point_app|

# マイグレーション

`$ make dry-migrate`でDDLを確認して、`$ make migrate`でマイグレーションする流れとなる

```sh
# マイグレーションする際に発光されるDDLを確認(実行はされない)
$ make dry-migrate
# マイグレーション適用
$ maae migrate
```

# 各ディレクトリの説明
詳しい説明は、各ディレクトリのREADME.mdに些細されているものもあります。

- `/handler`
  - ハンドラー層
  - クライアントのデータをバリデーション
  - クライアントにデータを返す
- `/service`
  - サービス層
  - ドメイン層、リポジトリ層を利用してユースケースを実現する
- `/domein`
  - ドメイン層
  - サービス間(ユースケース間)をまたがるロジックを記述
- `/repository`
  - リポジトリ層
  - DBやキャッシュサーバーにアクセスする
- `/entity`
  - DBエンティティをGoで医療するための構造体の定義
- `/_tools`
  - 各ツールを利用するためのディレクトリ
  - `/_tools/mysql/schema.sql`にテーブル作成のためのDDL定義
- `/config`
  - 各種環境変数などの設定値など
- `/utils`
  - ユーティリティパッケージ
