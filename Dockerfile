FROM golang:1.22-alpine

RUN apk add --no-cache curl

RUN curl -Lo /usr/local/bin/aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie-arm64
RUN chmod +x /usr/local/bin/aws-lambda-rie

WORKDIR /digitalvenue

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT [ "/usr/local/bin/aws-lambda-rie" ]
