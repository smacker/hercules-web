FROM golang:1.8-alpine3.6

# base deps
RUN apk --update upgrade && \
    apk add --no-cache make git curl ca-certificates bash \
    build-base libxml2-dev protobuf nodejs=6.10.3-r1 nodejs-npm && \
    npm install -g yarn

# install bblfsh go client deps
# otherwise hercules installation will fail
RUN go get -d -v gopkg.in/bblfsh/client-go.v2 && \
    cd $GOPATH/src/gopkg.in/bblfsh/client-go.v2 && \
    make dependencies

# install gogoproto (for hercules)
RUN go get github.com/gogo/protobuf/proto && \
    go get github.com/gogo/protobuf/jsonpb && \
    go get github.com/gogo/protobuf/protoc-gen-gogo && \
    go get github.com/gogo/protobuf/gogoproto

# install hercules (pinned version)
RUN mkdir -p $GOPATH/src/gopkg.in/src-d/hercules.v3 && \
    git clone -n https://github.com/src-d/hercules.git $GOPATH/src/gopkg.in/src-d/hercules.v3 && \
    git checkout 1f59ecd8043bc131efb40f2a221318c0dc144c00 && \
    cd $GOPATH/src/gopkg.in/src-d/hercules.v3 && \
    PATH=$PATH:$GOPATH/bin protoc --gogo_out=pb --proto_path=pb pb/pb.proto

ADD . /go/src/hercules-web
WORKDIR /go/src/hercules-web

# install deps
RUN go get github.com/jteeuwen/go-bindata/... && \
    go get && \
    yarn

# build everything
RUN yarn build && \
    go-bindata dist/... && \
    go build -o herculesweb .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/hercules-web/herculesweb .
CMD ["./herculesweb"]
