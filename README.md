# hercules-web

> web ui for [hercules](https://github.com/src-d/hercules)

## Build

```
docker build -t hercules-web .
docker run --rm -p 8080:8080 hercules-web
```

Open http://localhost:8080

## Development

### Frontend:

```bash
# install dependencies
yarn

# serve with hot reload at localhost:8081
yarn dev

# build dist files
yarn build
```

### Backend:

```bash
# install all deps (there are many of them, look at Dockerfile)
...

# server at localhost:8080
go run *.go
```
