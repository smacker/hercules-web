FROM golang:1.11-alpine3.8

# base deps
RUN apk --update upgrade && \
  apk add --no-cache make git curl ca-certificates bash \
  build-base libxml2-dev nodejs nodejs-npm && \
  npm install -g yarn

ADD . /hercules-web
WORKDIR /hercules-web

# install deps
RUN go get && \
  yarn

# build everything
RUN yarn build && \
  go build -o herculesweb .

FROM alpine:3.8
RUN apk --no-cache add ca-certificates libxml2 libgcc libstdc++
WORKDIR /root/
COPY --from=0 /hercules-web/dist ./dist
COPY --from=0 /hercules-web/herculesweb .
CMD ["./herculesweb"]
