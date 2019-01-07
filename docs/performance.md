# Performance

The injection of New Relic APM Metadata is implemented as a webhook using Kubernetes MutatingAdmissionWebhook, which is in fact an AdmissionWebhook. 

Admission webhooks run as part of the Kubernetes API Request Lifecycle.

> An admission controller is a piece of code that intercepts requests to the Kubernetes API server prior to persistence of the object, but after the request is authenticated and authorized. 
> 
> --- https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#what-are-they
 
One of our top priorities is to code it as lightweight as possible, making the webhook behave as non-intrusive and to perform almost at no cost.

Therefore, we ran the following benchmark and performance tests.

## Kubernetes Request lifecycle

### Internal Kubernetes HTTP API request flow
![](k8s-api-lifecycle.svg)

Any request to K8s API is affected by a 60 seconds (default value) timeout, including the time spent in all admission webhooks (and the whole lifecycle).
The default request timeout between the API and any admission webhook is 30 seconds as Kubernetes [source code](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apiserver/pkg/util/webhook/webhook.go#L36) shows.

### Caveats

* If a webhook does not give a response before 30 seconds, the K8s API will error because such timeout. This makes our hook to become a Highly Available service aiming to have 0 downtime.
* If the webhook does not give a response, the `failurePolicy` will not be used as it is only used when a response is given.

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
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    3000	    434999 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.420s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    2000	    604123 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.342s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    3000	    505401 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.647s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    2000	    642582 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.433s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    2000	    747675 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.662s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    2000	    837064 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.847s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    1000	   1005508 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.183s

$ go test -bench .                                                                                                                                                                                               
goos: darwin
goarch: amd64
pkg: go.datanerd.us/p/fryckbosch/k8s-env-inject
Benchmark_WebhookPerformance-8   	    3000	    447884 ns/op
PASS
ok  	go.datanerd.us/p/fryckbosch/k8s-env-inject	1.461s
```

## Benchmark having the Mutating Webhook in place

After running the code benchmark we realized that running benchmarks in real clusters are not important. As the service is simple and extremely fast, it should have a close to **zero** perceptive impact in pod creation in the cases where it would run.

We will use different tools to ensure performance is focussed in this service in the form of Golang benchmarks and a server side timeout to prevent perceptive interference in pod creation time.

