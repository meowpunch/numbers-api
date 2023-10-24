# Numbers API Service

## Requirements

`/numbers` endpoint accepts query parameter `u` about list of urls.
e.g.

```
http://127.0.0.1:8080/numbers?u=http://127.0.0.1:8090/primes&u=http://127.0.0.1:8090/fibo&u
=http://127.0.0.1:8090/rand
```

Then the service should retrieve `numbers` from urls.
And returns merged and sorted `numbers` in ascending order, which is that all number should be unique in the list.

### Additional Requirements

1. The endpoint needs to return the result as quickly as possible, but always within 500 milliseconds.
2. All URLs that were successfully retrieved within the given timeframe must influence the result of the endpoint.
3. If one URL takes longer to respond, keep loading it in the background and cache the response for future
   use. Pay very close attention to "CacheControl" instructions.
4. It is valid to return an empty list as a result only if all URLs returned errors or took too long to respond and no
   previous response is stored in the cache

## Assumptions

First of Additional requirements (AR-1) conflicts AR-2. If one URL takes exactly the given timeframe (500ms) to respond,
how can my service return the result within 500ms?
The service should do additional operation like merge and sort after retrieving numbers from the URL. It should take
over 500ms in the case.

I will achieve AR-2 fully and try to get closer to AR-1. That is, the endpoint might return over 500ms but 99 percentile
of request would take in 500ms.


## Test and Build


### Unit Test

```bash
go test ./...  -cover
```
