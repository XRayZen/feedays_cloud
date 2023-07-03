FROM golang:1.20 as builder
WORKDIR /write
COPY go.mod ./
# build
COPY . .
# ARM64向けにビルド
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build .
# ARMで動かすにはprovidedランタイムを使う
FROM public.ecr.aws/lambda/provided:al2
WORKDIR /write
COPY --from=builder /write/write /write
ENTRYPOINT [ "./write" ]