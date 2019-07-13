package bandithttp

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExperiment(t *testing.T) {
	mux := http.NewServeMux()
	s := &Experiment{
		Name: "foo a/b test",
		Arms: []*Arm{
			{
				Name: "A",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("I'm A"))
				}),
			},
			{
				Name: "B",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("I'm B"))
				}),
			},
		},
		Algorithm: EpsilonGreedy(0.1),
	}
	mux.Handle("/ab", s)
	server := httptest.NewServer(mux)
	client := server.Client()
	resp, err := client.Get(server.URL + "/ab")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

func TestEpsilonGreedy(t *testing.T) {
	mux := http.NewServeMux()
	s := &Experiment{
		Name: "foo a/b test",
		Arms: []*Arm{
			{
				Name: "A",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("I'm A"))
				}),
			},
			{
				Name: "B",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("I'm B"))
				}),
			},
		},
		Algorithm: EpsilonGreedy(0.1),
	}
	armRewardProb := map[string]float32{
		"A": 0.25,
		"B": 0.5,
	}
	mux.Handle("/ab", s)
	server := httptest.NewServer(mux)
	client := server.Client()
	armCount := map[string]int{
		"A": 0,
		"B": 0,
	}
	random := rand.New(rand.NewSource(42))
	for i := 0; i < 1000; i++ {
		resp, err := client.Get(server.URL + "/ab")
		if err != nil {
			t.Error(err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("unexpected status code: %d", resp.StatusCode)
		}
		usedArm := resp.Header.Get("Arm")
		armCount[usedArm] += 1
		prob := armRewardProb[usedArm]
		if random.Float32() < prob {
			body, err := json.Marshal(struct {
				Arm string `json:"arm"`
			}{usedArm})
			if err != nil {
				t.Fatal(err)
			}
			client.Post(server.URL+"/ab", "application/json", bytes.NewReader(body))
		}
	}
	if armCount["A"] > armCount["B"] {
		panic(armCount)
	}
}
