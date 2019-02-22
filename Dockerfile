FROM golang:1.11-alpine3.9

# base deps
RUN apk --update upgrade && \
  apk add --no-cache make ca-certificates git build-base nodejs yarn

ADD . /hercules-web
WORKDIR /hercules-web

# install deps
RUN go get && \
  yarn

# build everything
RUN yarn build && \
  go build -o herculesweb .

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /hercules-web/dist ./dist
COPY --from=0 /hercules-web/herculesweb .
CMD ["./herculesweb"]
