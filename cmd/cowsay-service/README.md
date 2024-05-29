# cowsay-service

Provides an http service that answers calls to following routes

## `/cowsays`

Returns the cowsay characters that are available

Example:

```bash
curl 'localhost:8080/cowsays'
```

## `/cowsay`

Example:

```bash
curl 'localhost:8080/cowsay' -d 'hello' 
```
