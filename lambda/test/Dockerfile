FROM golang:1.20 as builder
# 変数function_nameを作ってそこにtestを入れる
ARG function_name=test
WORKDIR /${function_name}
COPY go.mod go.sum ./
# build
COPY . .
# 必要なパッケージをインストール
RUN go get github.com/aws/aws-sdk-go/aws \
    && go get github.com/aws/aws-sdk-go/aws/session \
    && go get github.com/aws/aws-sdk-go/service/secretsmanager \
    && go get github.com/jinzhu/gorm/dialects/mysql@v1.9.16 \
    && go get github.com/mmcdole/gofeed
# ARM64向けにビルド
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build .
# ARMで動かすにはprovidedランタイムを使う
FROM public.ecr.aws/lambda/provided:al2
ARG function_name=test
WORKDIR /${function_name}
COPY --from=builder /${function_name}/${function_name} /${function_name}
ENTRYPOINT [ "./test" ]