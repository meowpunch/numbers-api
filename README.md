# Numbers API Service

## Requirements

- merge the numbers coming from all URLs, sort them in ascending order and make sure that each number appears only once
  in the result (usecase/merger/merger_test.go)
- an URL is not valid or does not return the correct result simply ignore it. (usecase/fetcher/fetcher_test.go)

### Additional Requirements

1. The endpoint needs to return the result as quickly as possible, but always within 500 milliseconds. (
   usecase/get_numbers_test.go)
2. All URLs that were successfully retrieved within the given timeframe must influence the result of the endpoint. (
   usecase/get_numbers_test.go)
3. If one URL takes longer to respond, keep loading it in the background and cache the response for future
   use. Pay very close attention to "CacheControl" instructions. (usecase/fetcher/fetcher_test.go)
4. It is valid to return an empty list as a result only if all URLs returned errors or took too long to respond and no
   previous response is stored in the cache (usecase/get_numbers_test.go)

### Assumptions

First of Additional requirements (AR-1) conflicts AR-2. If one URL takes exactly the given timeframe (500ms) to respond,
how can my service return the result within 500ms?
The service should do additional operation like merge and sort after retrieving numbers from the URL. It should take
over 500ms in the case.

I will achieve AR-2 fully and try to get closer to AR-1. That is, the endpoint would return over 500ms but 99 percentile
of request would take in 500ms. (integration/load_test.go)

## Get Started

### Test

Unit tests

```bash
go test -v ./... -tags='!integration'
```

Integration test (make sure if 8080 and 8090 ports are not occupied)

```bash
go test -v ./integration/load_test.go
```

### Run locally

run test server and numbers api server

```bash
docker run --detach --publish 8090:8090 emanuelschmoczer/coding-challenge-test-server:latest
go run main.go
```

request to the endpoint

```bash
curl http://localhost:8080/numbers?u=http://127.0.0.1:8090/primes&u=http://127.0.0.1:8090/fibo&u=http://127.0.0.1:8090/rand&u=http://127.0.0.1:8090/rand
```

latency of the request

```bash
curl -o /dev/null -s -w "Connect: %{time_connect} TTFB: %{time_starttransfer} Total time: %{time_total} \n"  http://localhost:8080/numbers?u=http://127.0.0.1:8090/primes&u=http://127.0.0.1:8090/fibo&u=http://127.0.0.1:8090/rand&u=http://127.0.0.1:8090/rand
```

## Furthermore

I'm convinced that the project has robust test to achieve correctness.
However, I cannot say it's production ready because there are missing parts.
Especially, I approximately guarantee latency within 500ms. In that case, monitoring system is helpful to check
performance. 