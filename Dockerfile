FROM golang:latest

#ENV GO111MODULE=on

WORKDIR /web/gin-photo-storage

#COPY go.mod go.sum ./

#RUN apk add git && go mod download

COPY . .

RUN go build

CMD ["./gin-photo-storage"]

ENTRYPOINT ["./bootstrap.sh"]
