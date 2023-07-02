FROM golang:1.20 as builder
WORKDIR /heavy
COPY go.mod go.sum ./
# build
COPY . .
# ARM64向けにビルド
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build .
# ARMで動かすにはprovidedランタイムを使う
FROM public.ecr.aws/lambda/provided:al2
WORKDIR /heavy
COPY --from=builder /heavy/heavy /heavy
ENTRYPOINT [ "./heavy" ]