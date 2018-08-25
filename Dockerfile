FROM golang:1.10-alpine3.8

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
  git checkout 1bdd0b79356e89da1a6d1558f15be87d4f722e87 && \
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
