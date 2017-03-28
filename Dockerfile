FROM golang:1.7

ENV INDEX /data/cardindex

ADD . /go/src/github.com/arjandepooter/discord-epic-cardbot
WORKDIR /go/src/github.com/arjandepooter/discord-epic-cardbot

RUN go get
RUN go install github.com/arjandepooter/discord-epic-cardbot

VOLUME /data

ENTRYPOINT ["/go/bin/discord-epic-cardbot"]
