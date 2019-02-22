# hercules-web

> web ui for [hercules](https://github.com/src-d/hercules)

![Project Burndown](docs/project-burndown.png?raw=true)

## Run

```
docker run --rm -p 8080:8080 smacker/hercules-web
```

Open http://localhost:8080

## Development

### Frontend:

```bash
# install dependencies
yarn

# serve with hot reload at localhost:8081
yarn serve

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
