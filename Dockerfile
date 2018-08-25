FROM golang:1.8-alpine3.6

# base deps
RUN apk --update upgrade && \
  apk add --no-cache make git curl ca-certificates bash \
  build-base libxml2-dev protobuf nodejs nodejs-npm && \
  npm install -g yarn

# install gogoproto (for hercules)
RUN go get github.com/gogo/protobuf/proto && \
  go get github.com/gogo/protobuf/jsonpb && \
  go get github.com/gogo/protobuf/protoc-gen-gogo && \
  go get github.com/gogo/protobuf/gogoproto

# install hercules (pinned version)
RUN mkdir -p $GOPATH/src/gopkg.in/src-d/hercules.v3 && \
  git clone -n https://github.com/smacker/hercules.git $GOPATH/src/gopkg.in/src-d/hercules.v3 && \
  cd $GOPATH/src/gopkg.in/src-d/hercules.v3 && \
  git checkout 4e3b0054e9a88bce5b3941cc568b7af1b6a71c2b && \
  PATH=$PATH:$GOPATH/bin protoc --gogo_out=pb --proto_path=pb pb/pb.proto

ADD . /go/src/hercules-web
WORKDIR /go/src/hercules-web

# install deps
RUN go get && \
  yarn

# build everything
RUN yarn build && \
  go build -o herculesweb .

FROM alpine:latest
RUN apk --no-cache add ca-certificates libxml2 libgcc libstdc++
WORKDIR /root/
COPY --from=0 /go/src/hercules-web/dist ./dist
COPY --from=0 /go/src/hercules-web/herculesweb .
CMD ["./herculesweb"]
