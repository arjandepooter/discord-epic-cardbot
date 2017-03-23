FROM golang:1.7

ADD . /go/src/github.com/arjandepooter/discord-epic-cardbot
WORKDIR /go/src/github.com/arjandepooter/discord-epic-cardbot

RUN go get
RUN go install github.com/arjandepooter/discord-epic-cardbot

ENTRYPOINT ["/go/bin/discord-epic-cardbot"]
