FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

ARG GO_CMD
ENV GO_CMD=${GO_CMD}

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/${GO_CMD} github.com/DarioRoman01/delfos-chat/cmd/${GO_CMD}


FROM golang:1.19-alpine
ARG GO_CMD
ENV GO_CMD=${GO_CMD}

COPY --from=builder /app/bin/${GO_CMD} /usr/local/bin/${GO_CMD}
CMD ${GO_CMD}