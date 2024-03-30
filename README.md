# lamoda-test

## Start app

```
make up
```

## Run tests

```
make test
```

## API
### Check remaining goods in store

#### Request:
```
curl -H 'Content-Type: application/json' -X POST -d '{"store": 1}' http://localhost:8000/check_goods
```

#### Response:
```
{
    "goods":[
        {"uuid":"1720f137-4e06-427a-aa0c-6b22c35eecc6","name":"goods1","size":"50x50x10","amount":100},
        {"uuid":"399861f6-6f57-413d-97cf-c73b3ab09de1","name":"goods2","size":"10x10x10","amount":200},
    ]
}
```

### Reserve goods in store

#### Request:
```
curl -H 'Content-Type: application/json' -X POST -d '{"goods": [
    {"uuid": "1720f137-4e06-427a-aa0c-6b22c35eecc6", "amount": 10}, 
    {"uuid": "399861f6-6f57-413d-97cf-c73b3ab09de1", "amount": 10}]}' http://localhost:8000/reserve_goods
```

#### Response:
```
{
    "result":[
        {"uuid":"1720f137-4e06-427a-aa0c-6b22c35eecc6","status":"OK"},
        {"uuid":"399861f6-6f57-413d-97cf-c73b3ab09de1","status":"OK"}
    ]
}
```

### Free goods in store

#### Request:
```
curl -H 'Content-Type: application/json' -X POST -d '{"goods": [
    {"uuid": "1720f137-4e06-427a-aa0c-6b22c35eecc6", "amount": 10}, 
    {"uuid": "399861f6-6f57-413d-97cf-c73b3ab09de1", "amount": 10}]}' http://localhost:8000/free_goods
```

#### Response:
```
{
    "result":[
        {"uuid":"1720f137-4e06-427a-aa0c-6b22c35eecc6","status":"OK"},
        {"uuid":"399861f6-6f57-413d-97cf-c73b3ab09de1","status":"OK"}
    ]
}
```
