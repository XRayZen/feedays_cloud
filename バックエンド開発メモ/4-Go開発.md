
# GoでLambda開発
- [x] Goの概要と文法を調べメモを取る
- [x] GoでLambdaを開発する方法を調べメモを取る
- [x] GoでLambdaを開発する

## まずは最小構成でLambdaを実装・デプロイしてHelloWorldを試す
- GoとDockerfileを実装完了したのでinfraをECR+Lambdaの最小構成で実装してデプロイする
- ECRにプッシュするのはインフラがデプロイした後にする
- 大半のコードはTestCICDを流用できる
- AWS ECRのイメージをリスト
```bash
aws ecr list-images --repository-name feedays-cloud-repo
```
- ECRからイメージを削除する
```bash
aws ecr batch-delete-image --repository-name feedays-cloud-repo --image-ids imageTag=world
```
- Lambda関数コードをDockerイメージにビルドする
```bash
docker build -t ＜AWS Account ID＞.dkr.ecr.us-east-1.amazonaws.com/feedays-cloud-repo:read ./lambda/read
```
- AWS ECRにログインしてDockerログインする
```bash
aws ecr get-login-password | docker login --username AWS --password-stdin ＜AWS Account ID＞.dkr.ecr.us-east-1.amazonaws.com/feedays-cloud-repo:world
```
- ECRにDockerイメージをプッシュする
```bash
docker push ＜AWS Account ID＞.dkr.ecr.us-east-1.amazonaws.com/feedays-cloud-repo:world
```
- Lambda関数を更新する
```bash
aws lambda update-function-code --function-name hello-world --image-uri ＜AWS Account ID＞.dkr.ecr.us-east-1.amazonaws.com/feedays-cloud-repo:world
```

- ECRイメージ一覧を取得 
```bash
aws ecr describe-images --repository-name feedays-cloud-repo
```
## エラーが解決してHello-Worldを実行できたので、次はGoのテストを調べメモして実装する
- 2023年6月8日概ね完了した
## Goにテストをし易くするためにDIを導入する
- Lambda関数は他のサービスを呼び出すため、テストをし易くするためにDIを導入する
- 呼び出すAWSサービスをリポジトリにまとめて抽象化してテスト時本番時に切り替える
# OpenAPIの定義をコードを生成する
- まずはAPIをパースする
  - OpenAPIの定義をコードを生成する
- プロジェクトルートで以下のコマンドを実行する
```bash
oapi-codegen -config OpenAPI/api_gen_config.yaml OpenAPI/backend_api.yaml > OpenAPI/backend_api.gen.go
```
# まずはread系関数を設計どおりに実装する
  - リクエストを生成してコードに通して変換する
  - インポートできないエラーがあったがGoはインポートするパッケージ名とディレクトリ名は一致していないといけない
- リクエストタイプに応じて処理する
  - ExploreCategory
    - 国別に応じてExploreCategoryを返すための処理を実装する
    - ユーザーIDをキーにしてユーザー情報を取得して国を取得して国別のExploreCategoryを返す
  - Ranking
    - 国別に応じてランキングを返すための処理を実装する
    - ユーザーIDをキーにしてユーザー情報を取得して国を取得して国別のランキングを返す
  - ユーザー情報を取得
    - APIリクエスト制限・ユーザー情報を取得してクライアント側で色々と処理をする
    - ユーザー情報の定義をここでコードに起こすと大部分は使われない
      - 使うのはSQLだから取得する情報を制限すれば何も問題ない
- 帰ってきた情報をStructに変換してレスポンスに入れる時にjsonに変換して返す
- DIをする為にGoでのDIを調べる
- Responseを返そうとしたがレスポンス用のコードが生成されていない
  - API定義を修正したら生成された
- テストも正常系と異常系を実装してテストをパスしたから概ね実装完了

## Read系を概ね八割実装したので、次はWrite系を実装する
- まずはブランチを切る
- ブランチ名は`feature/Impl-lambda-write`とする
  - `git checkout -b feature/Impl-lambda-write`
- ユーザー情報を取得を実装する
  - ユーザー情報を取得する時にユーザーIDをキーにして各テーブルから取得する
  - ユーザー情報にはUI/購読サイトなどの設定情報も入っているのでここで実装する
- 定義するデータ型には全てJsonをタグをつけてクライアントのJsonタグと合わせて、変換に齟齬が出ないようにする
- [x] api_functionsを実装する
- [x] `ユーザーを登録`を実装
  - 受け取ったjsonをユーザー設定に変換してDBに登録する
  - API応答を作る
  - ユーザー登録をイベントとしてアクティビティレコードに書き込む
- [x] ユーザー識別情報も新しく定義するアクテビティレコードに入れて実装する
  - [x] APIにアクティビティレコード読み書きを実装する
- アクティビティレコードの定義
  - アクティビティID
  - アクティビティタイプ
    - API機能に応じる
  - アクティビティ日時
  - ユーザーID
  - アクセス環境タイプ
    - Webかモバイルかデスクトップか
  - アクセスIP
- アクティビティレコードを作る処理を実装する
  - Goで現在日時を取得する
### テストコードを実装する
-  UserAccessIdentInfoのjson復号化に失敗する
   -  json: cannot unmarshal object into Go value of type string
   -  Unmarshalのコードに問題があって修正したら治った
- 一応、テストをパスしたから概ね実装完了
  - 残りはHeavy関数を実装する
  - そして、DB関連を実装する

## Docker-composeでGo開発環境を構築する
- アピール材料になる為、ここでこの作業をしておく必要がある
- 方法を調べてメモをする

## 重い処理をするLambda関数を実装する
- 関数名は`heavy`とする
- AWSリソースを多めに振る
- 検索を実装する
  - キーワード検索ならDBからキーワードを含む記事を取得する
  - サイト検索ならDBからサイトを取得する
    - サイト検索はサイト名をキーにしてDBからサイトを取得する
  - URLならDBからURLを取得する
  - 未登録なら登録処理をする
- 他のGoモジュールから型を参照できるのでreadで集約して定義する
  - DBRepoも集約して定義する
  - もし、間違えていたらDRY原則に反するがDBRepoは別々で実装するしかない
  - それでもダメなら、保守性が悪くなるがData定義を別々にするしかない
  - DockerfIleでのビルド上分割しての実装が得策
    - 今は集約しておいてLambdaにデプロイする時に分割する
    - その際は全てのデータ型・DBRepoをコピーしておく
      - それをやるとコードが重複するが、それは仕方ない
      - 管理が手間になる為、要約してから個別にフォルダごとコピーする
- 検索を実装する
  - 今の所URL検索とキーワード検索（記事）のみだが、今後はサイト名検索も実装する
  - その為にリクエストに検索タイプを追加する
  - 後はテストケースを網羅してテストをパスさせる
- 購読を実装する
  - テストをパスして概ね実装完了

## 定期的に実行するLambda関数を実装する
- イベントで実行する
- サイトテーブルから読み込んで記事を更新する
  - それにより購読サイトのフィードの鮮度を維持する
  - 並列処理で実装する
  - テストをパスして概ね実装完了
- モック用にリードアクテビティを作る
  - リードアクテビティに足りない要素があるので追加する
    - 実装できたのでクライアント側もやる

!!! info ランキング関連は実装が複雑化するので、後回しにする
    時間が圧しているので、リードアクテビティは集めるがランキング機能はリリース後に実装する

