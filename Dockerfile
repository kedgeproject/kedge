FROM golang:1.8

WORKDIR /go/src/github.com/kedgeproject/kedge
COPY . .

RUN go install

ENTRYPOINT ["kedge"]
