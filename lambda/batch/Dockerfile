FROM golang:1.20 as builder
WORKDIR /batch
COPY go.mod go.sum ./
# build
COPY . .
# 必要なパッケージをインストール
RUN go get github.com/aws/aws-sdk-go/aws \
    && go get github.com/aws/aws-sdk-go/aws/session \
    && go get github.com/aws/aws-sdk-go/service/secretsmanager \
    && go get github.com/aws/aws-secretsmanager-caching-go/secretcache \
    && go get gorm.io/gorm \
    && go get gorm.io/driver/mysql \
    && go get gorm.io/driver/sqlite \
    && go get github.com/mmcdole/gofeed

# ARM64向けにビルド
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build .
# ARMで動かすにはprovidedランタイムを使う
FROM public.ecr.aws/lambda/provided:al2
WORKDIR /batch
COPY --from=builder /batch/batch /batch
ENTRYPOINT [ "./batch" ]