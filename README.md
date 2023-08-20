# FeedaysのバックエンドAPI


以下使用技術一覧やコンセプト

## 開発環境
Dockerfileで開発コンテナを記述してVSCodeのRemote Containerで開発環境を構築しています。




curl -X POST -H "Content-Type: application/json" -d '{"userId":"test","requestType":"ServiceInitialize","requestArgumentJson1":"","requestArgumentJson2":""}' https://bkq8lpslz8.execute-api.us-east-1.amazonaws.com/develop/user
