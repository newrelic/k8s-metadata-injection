# Performance

The injection of New Relic APM Metadata is implemented as a webhook using Kubernetes MutatingAdmissionWebhook, which is in fact an AdmissionWebhook. 

Admission webhooks run as part of the Kubernetes API Request Lifecycle.

> An admission controller is a piece of code that intercepts requests to the Kubernetes API server prior to persistence of the object, but after the request is authenticated and authorized. 
> 
> --- https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#what-are-they
 
One of our top priorities is to code it as lightweight as possible, making the webhook behave as non-intrusive and to perform almost at no cost.

Therefore, we ran the following benchmark and performance tests.

## Kubernetes Request lifecycle

Please refer to the [Request internal lifecycle](lifecycle.md) documentation.

## Benchmark of the webhook code

The MutatingAdmissionWebhook code has been tested using Golang Benchmarks.

The slowest run was on average 1005508 ns/op. The fastest on average was 434999 ns/op. The "average of the average" was 653154 ns/op. All these values are 1 millisecond or less. On a real word situation it will be slower due to different size of each pod's creation payload and to the TLS overhead.

The code of the benchmark can be found [here](./webhook_test.go).

### Results

* These tests were ran on a 2018 Macbook Pro with a core i7 2.7 GHz and 16 GB of memory.
* They were recorded on Jan 4th, 2019.
* Latest commit SHA was `ab4b0de131c4e2ea51089e1441d5e571baa3803e`


```
$ go test -bench .                                                                                                                                                                                               
                                  
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    3000	    434999 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.420s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    2000	    604123 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.342s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    3000	    505401 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.647s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    2000	    642582 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.433s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    2000	    747675 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.662s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    2000	    837064 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.847s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    1000	   1005508 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.183s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: github.com/newrelic/k8s-metadata-injection
Benchmark_WebhookPerformance-8   	    3000	    447884 ns/op
PASS
ok  	github.com/newrelic/k8s-metadata-injection	1.461s
```

## Benchmark having the Mutating Webhook in place

After running the code benchmark we realized that running benchmarks in real clusters are not important. As the service is simple and extremely fast, it should have a close to **zero** perceptive impact in pod creation in the cases where it would run.

We will use different tools to ensure performance is focused in this service in the form of Golang benchmarks and a server side timeout to prevent perceptive interference in pod creation time.

