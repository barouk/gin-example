#FROM golang:1.18.1-alpine3.15
FROM  2c53c1429c40


WORKDIR  /nic-core-handles
COPY --from=nic-core-common:latest /nic-core-common/ /nic-core-common/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "main.go"]
