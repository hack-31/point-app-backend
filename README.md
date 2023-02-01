# point-app-backend
ハッカソン2022で作成するポイントアプリのバックエンドプロジェクト

# 開発手順

```sh
# 環境手順
$ git clone https://github.com/hack-31/point-app-backend.git
$ cd ./point-app-backend

# JWTに必要なキー生成
$ make create-key
# ビルド、コンテナ起動
$ make build-up
# マイグレーション適用
$ maae migrate
# サーバー起動
$ make serve
```

詳しくは、makeコマンドを実行

サーバー起動したら、ホスト側で以下のURLでアクセス可能

- API http://localhost:8081
- swagger-ui http://localhost:80
- adminer http://localhost:8082
- aws http://localhost:4566

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

# 初期データ挿入

`./_tool/mysql/seed.sql`に初期データを挿入するコマンドを記述
AWS cloud9の[point-app-dev](https://ap-northeast-1.console.aws.amazon.com/cloud9/home/environments/9e3ee1e0dda0408b80b541ecd88be4da?permissions=owner)で以下のコマンドで挿入

```sh
$ make seed
```

# 各ディレクトリの説明
詳しい説明は、各ディレクトリのREADME.mdに些細されているものもあります。
- `/auth`
  - JWT認証におけるトークン作成や検証を行う処理が書かれたパッケージ
- `/router`
  - ルーティング処理を書きます
- `/docs`
  - ドキュメントが配置
  - swaggerなどを修正する際は`/docs/openapi.yml`で行う
- `/handler`
  - ハンドラー層
  - クライアントのデータをバリデーション
  - クライアントにデータを返す
- `/service`
  - サービス層
  - ドメイン層、リポジトリ層を利用してユースケースを実現する
- `/domain`
  - ドメイン層
  - サービス間(ユースケース間)をまたがるロジックを記述
- `/repository`
  - リポジトリ層
  - DBやキャッシュサーバーにアクセスする
- `/_tools`
  - 各ツールを利用するためのディレクトリ
  - `/_tools/mysql/schema.sql`にテーブル作成のためのDDL定義
- `/config`
  - 各種環境変数などの設定値など
- `/constant`
  - サービス全体として定数などを定義
- `/utils`
  - ユーティリティパッケージ

# デプロイ
GitHub Action, Codepipelineによる自動デプロイになります。mainブランチにマージされると自動でデプロイされます。
GitHub Actionでは、Docker Imageを作成、ECRに登録を行い、Codepipelineでは、ECRからImageを取得し、ECSにデプロイします。

## 手動デプロイ方法
AWS cloud9の[point-app-dev](https://ap-northeast-1.console.aws.amazon.com/cloud9/home/environments/9e3ee1e0dda0408b80b541ecd88be4da?permissions=owner)にて以下のコマンド実行

```sh
$ pwd
/home/ec2-user/environment/point-app-backend
# 値の設定
$ AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query 'Account' --output text)
$ IMAGE_TAG=$(git rev-parse HEAD)
$ ECR_REGISTRY=${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com
# ビルド
$ docker image build \
  -t ${ECR_REGISTRY}/point-app-backend:latest \
  -t ${ECR_REGISTRY}/point-app-backend:${IMAGE_TAG}  \
  --target deploy ./ --no-cache
# ECRログイン
$ aws ecr --region ap-northeast-1 get-login-password | docker login --username AWS --password-stdin ${ECR_REGISTRY}/point-app-backend
# ECR push
$ docker image push -a ${ECR_REGISTRY}/point-app-backend
```

ECRにデプロイされると、latestタグの更新がトリガーとなり、codepipelineによってECSにデプロイが走ります。
