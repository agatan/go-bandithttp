## bandithttp: Multi-Armed Bandit Problem on HTTP servers.

NOTE: This is just an experimental project.

`bandithttp` provides algorithms and utilities for multi-armed bandit problem with `http` package.


## Example


```go
urlA, _ := url.Parse("https://foo.example.com")
urlB, _ := url.Parse("https://bar.example.com")
exp := &bandithttp.Experiment{
	Name:      "TheExperiment",
	Algorithm: bandithttp.EpsilonGreedy(0.1),
	Arms:      []*bandithttp.Arm{
		{Name: "foo", Handler: httputil.NewSingleHostReverseProxy(urlA)},
		{Name: "bar", Handler: httputil.NewSingleHostReverseProxy(urlB)},
	},
}
http.ListenAndServe(":8080", exp)
```
