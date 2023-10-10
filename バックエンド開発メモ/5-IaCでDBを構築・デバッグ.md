- [AWS RDSを調査](#aws-rdsを調査)
- [AWS側でRDSを構成する](#aws側でrdsを構成する)
- [TerraformでRDS Proxyを構成する](#terraformでrds-proxyを構成する)
- [Terragruntで適用する](#terragruntで適用する)
  - [Test用のLambda関数を作って検証する](#test用のlambda関数を作って検証する)
    - [最後にLambdaのシークレット取得周りで無駄なIAMがあるのでそれを取り除いてテストする](#最後にlambdaのシークレット取得周りで無駄なiamがあるのでそれを取り除いてテストする)
    - [Lambda関数更新手順](#lambda関数更新手順)

# AWS RDSを調査
- [x] AWSのリレーショナルDBサービスを調べる
  - [x] RDSProxyを調べる
  - [x] ベストプラクティスを調べる
# AWS側でRDSを構成する
- RDSProxy＋RDSで構成する
- DBエンジンにMySQLを選択する
  - 使うバージョンは8.0.31
- 使うインスタンス
  - 無料枠では`db.t3.micro`
    - 開発環境ではこれで良い
  - 本番環境でも、`db.t3.micro`を引き続き使う予定
  - 就職活動時はアップグレードして`db.t3.medium`を使う
    - これでRDS単体予測コストが月51.94USD
- ストレージ
  - 割り当ては20GiB
  - 最大ストレージは22GiB
    - これは開発・本番環境で同じ
    - サービスをロンチするとしても最小コストに抑えることを最優先
    - 就職活動時でも予兆がない場合は20GiBで良い
- パラメータグループ
  - パラメータグループについては調べる
  - `default.mysql8.0`を使う
- マルチAZではなくシングルAZ
- バックアップ/削除設定・DBサブネットグループの定義を完了
- データベースインスタンスの定義を完了
# TerraformでRDS Proxyを構成する
- RDSに接続する時にSecretsManagerを使う必要があるのか
- RDSProxyのスタート方法を調べてメモする
- ネットワークの設定をする
  - プライベートサブネットに配置するためにインプットからの変数をセットする
    - Lambdaなどと同じVPCに配置する
  - セキュリティグループをセットする
    - [これ](https://github.com/terraform-aws-modules/terraform-aws-rds/blob/master/examples/complete-mysql/main.tf)を参考にセキュリティグループをモジュール化してセットする
  - モジュールVPCにデータベース用のサブネットグループが用意されているからそれを使用してセットする
- RDSProxyも[これを](https://github.com/cloudposse/terraform-aws-rds-db-proxy/blob/main/examples/complete/main.tf)参考にしてモジュールを使って定義してLambdaと接続する
- GoのLambdaでは`os.Getenv("dbUser")`で取得できる
  - GOメモに書いておく
- LambdaがRDS Proxyにアクセスするには環境変数にエンドポイントをセットしてIAMロールにもRDS PRoxyのARNをセットして権限を付与する必要がある
  - `aws_iam_policy_attachment`を使って権限を付与する
    - これは複数定義できるからマネージドポリシーと合わせて使える
  - [これ](https://github.com/aws-samples/serverless-patterns/blob/main/apigw-http-api-lambda-rds-proxy-terraform/main.tf)をLambdaとRDSProxy接続の参考にする
- LambdaがVPC内のリソースとインターネットに接続するにはプライベートサブネットに配置してパブリックにNATゲートウェイを配置する必要がある
  - NATゲートウェイは無料枠の対象外で35USD/月
  - RDSは無料枠の範囲内で運用してDNS＋Route53を追加するとどうなるのか
    - これの段階はWebフロントエンド実装時に検証する
  - 今回はLambdaはRDS Proxyと同じサブネット内にLambdaを配置する
- RDSとVPCの設定が完了
- 次はLambdaの設定をする
  - RDSとインターネットに接続するからプライベートサブネットに配置する
  - 分かっていないのはモジュールVPCにはインターネットゲートウェイを設定しなくても良いのかという問題
    - モジュールにはIGWやACLがついてるから自動的に設定されているから大丈夫
# Terragruntで適用する
- [x] AWSCLIに認証情報をセットする
  - [x] コンテナを開発コンテナでできないからローカルでやる

## Test用のLambda関数を作って検証する
- これはLmabda＋RDSでDBのCRUDが動くのかを検証する為
  - [x] シークレット読み込みの検証が確認したら`secretcahche`に切り替えて検証する
    - AWS公式曰く呼び出しにコストがかかるからキャッシュ経由での呼び出しを推奨している
  - [x] 起動したらDBに接続してテーブルを作成してデータをCRUDする
- Lambda関数から、SecretsManagerに GetSecretValue を実施すると timeout するエラーが発生した
  - Lambdaのセキュリティグループにポート443のアウトバンドを許可する必要があるとのこと
  - 設定したら`GetSecretValueInput`のエラーは解消した
- だが、`GetSecretValue`でタイムアウトするようになった
  - セキュリティグループのアウトバンドを無制限にしたら解消した
  - [x] シークレットを取得できたのでシークレットキャッシュ経由で取得するようにする
- RDSコネクトで接続認証エラーが発生した
  - `{"errorMessage":"driver: bad connection","errorType":"errorString"}`
  - dbNameとdbPortと&&を追加してもダメだった
  - RDSのセキュリティグループを検証してもダメだった
  - コンソールでRDSを確認したらデータベース識別子にDB名-rdsとついていた
    - DB名は間違えていなかった
  - RDSモジュールで`target_db_instance = true`を忘れていた為かコンソールではプロクシーにターゲットDBが存在しなかった
  - TLSで接続出来ないので無効化した
  - `Database Access denied for user 'admin'@'10.0.3.6`でアクセスを拒否された
    - 通ってるという事はセキュリティーグループの設定は問題ない
    - RDSのDBにユーザー登録をしていないからアクセスを拒否されている
    - その為にRDSのDBにユーザー登録をするEC2インスタンスを作成する必要がある
      - モジュールを使ってコード量を極力減らす
      - テストが完了して用が終わったらすぐに削除する
- [x] テスト用にEC2インスタンスを作ってDBにユーザー登録する
  - テスト用のEC2を作成する為に色々と変更したから一旦削除して出直す
  - 作成コマンドはstage/develop/test_ec2で実行する
    - `terragrunt apply`
  - 再作成時にシークレットのエラーが出た
    - `You can't create this secret because a secret with this name is already scheduled for deletion.`
    - この[クラスメソッド記事](https://dev.classmethod.jp/articles/secrets-manager-error-recovery-window/)を参考にしてAWS CLIで強制的に削除した
- 無事にmysqlにアクセスできたのでエラーをまとめる
  - 原因はただ単にパスワードが間違っていたから
  - これでRDSProxyを使ってLambdaからRDSにアクセスできる
  - lambda_testの関数コンテナをデプロイして起動する
- コンソールでRDS側のシークレットを調べたら、RDSProxyのとパスワードが違っていた
  - これがRDSProxyの認証は通ってもRDSにアクセスできなかった原因か
  - 設定に誤りがあるか検証
  - おそらくはrandom_passwordのresultをlocalsを挟まずに直接参照したのが原因
    - [参考元](https://github.com/terraform-aws-modules/terraform-aws-rds-proxy/blob/master/examples/mysql-iam-instance/main.tf)はちゃんと挟んでいたが仕様をあまり調べずに安易に使って時間を無駄にしてしまった
  - これをlocalsを挟んで参照するように修正した
  - 今RDSでapplyしているがlambdaでもする必要がある
  - 変更を完全に反映する為にRDSを一旦削除して再作成する
  - [x] そして、RDSとProxyのシークレットが一致しているのか調べる
    - 調べたら、また違っていたので先にシークレットをAWSCLIで強制的に削除してからデストロイしてアプライする
    - シークレット強制削除コマンド
      - `aws secretsmanager delete-secret --secret-id admin --force-delete-without-recovery`
    - 完璧にしてもローテーション設定などが違うから原因はこれではないだろう
  - RDSデプロイが完了したらLambdaのデプロイをする
- 結局、`Database Access denied for user`にはまり続ける
  - RDSProxyとの接続認証には成功している
  - そこから先のRDSとの接続認証が拒否されている
  - 次にRDSのSGを確認しているが特に問題が見つからない
    - [x] rds_proxyのSGにデータベースサブネットからのアクセスを許可するインバウンドを設定して試す
      - ダメだった
  - Proxyから先にアクセス出来ない
  - RDSのパスワードが違っている事が原因か
    - RDSのパスワードを無効化して再設定する
- RDSにマスターパスワードを管理させないようにしたらテストが成功した
  - エラー解決はマークダウンに記載して後でZennにまとめる
- EC2も作成する必要が無いのでEC2のドキュメントにコードを移動する


- [x] モジュールRDSProxyのエンドポイントでリードオンリーのはずなのにRWできるから制限できる方法を調べる
  - `tagete_role= "READ_ONLY"`を追加

### 最後にLambdaのシークレット取得周りで無駄なIAMがあるのでそれを取り除いてテストする
  - 具体的にはシークレットマネージャーへのアクセス権限を付与する部分
  - すでに付与したマネージドポリシーで十分なので不要
### Lambda関数更新手順
- ログイン
```bash
aws ecr get-login-password | docker login --username AWS --password-stdin ＜AWS AccountID＞.dkr.ecr.us-east-1.amazonaws.com/feedays_cloud-repo:test
```
- 最初にECRにある既存のイメージを削除する
```bash
aws ecr batch-delete-image --repository-name feedays_cloud-repo --image-ids imageTag=test
```
- Lambda関数コードをDockerイメージにビルドする
```bash
docker build -t ＜AWS AccountID＞.dkr.ecr.us-east-1.amazonaws.com/feedays_cloud-repo:test ./lambda/test
```
- ECRにプッシュする
```bash
docker push ＜AWS AccountID＞.dkr.ecr.us-east-1.amazonaws.com/feedays_cloud-repo:test
```
- Lambda関数を更新する
```bash
aws lambda update-function-code --function-name feedays-cloud-test --image-uri ＜AWS AccountID＞.dkr.ecr.us-east-1.amazonaws.com/feedays_cloud-repo:test
```
