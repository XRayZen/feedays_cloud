# FeedaysのバックエンドAPI
このプロジェクトは開発しているRSSリーダーサービスのサーバーレスなバックエンドAPI部分です
サーバーサイドの言語にGoを採用しており、SQLはGORMを使用しています。
インフラはAWSを採用しており、APIGW＋Lambda+RDSProxy＋RDS
IaCはTerraform＋Terragruntを採用しており、開発・運用を見据えて開発環境・本番環境を切り替えることが出来ます
# クラウド構成
APIGW->Lambda(Go)->RDSProxy->RDS(MySQL)

## 以下使用技術一覧やコンセプト
- ラムダ動作言語: Go
- IaC: Terraform+Terragrunt(IaCの保守・運用性を効率的にするツール)

!!! info Github ActionsでCI/CDはまだ未実装です

## 開発環境
Dockerfileで開発コンテナを記述してVSCodeのRemote Containerで開発環境を構築しています。

