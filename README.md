# Api Service

## Make
### Build
```bash
make -f ./scripts/Makefile build
```

### Run
```bash
make -f ./scripts/Makefile run
```

### Clean
```bash
make -f ./scripts/Makefile clean
```

## Docker

### Build
```bash
docker build -t api-service -f ./.docker/Dockerfile .
```

### Run
```bash
docker run -p 8080:8080 api-service
```

requests:
```
curl -X POST \                
  http://localhost:8083/api/v1/task \
  -H 'Content-Type: application/json' \
  -d '{"name": "Новая задача", "difficulty": 1}'
curl -X GET \ 
  http://localhost:8083/api/v1/task/123
```
correct difficulty [0, 10]