gRPC powered pair of microservices that process (parse&save) CSV data.

Built and run with:
```sh
go build && ./cservices
```

Run a docker container:
```sh
docker built -t cservices .
docker run --publish 7777:7777 --name ccservices
```
