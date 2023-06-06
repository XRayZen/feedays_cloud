FROM golang:1.20 as builder
WORKDIR /read
COPY go.mod go.sum ./
# build
COPY . .
# ARM64向けにビルド
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o read
# 実行用イメージ
FROM public.ecr.aws/lambda/provided:al2
COPY --from=builder /read/read /main
CMD ["/main"]