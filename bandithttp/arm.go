package bandithttp

import "net/http"

type Arm struct {
	Name    string
	Handler http.Handler
	trial   int
	reward  int
}

func (a *Arm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Arm", a.Name)
	a.Handler.ServeHTTP(w, r)
}

func (a *Arm) Trial() int {
	return a.trial
}

func (a *Arm) Reward() int {
	return a.reward
}

func (a *Arm) ExpectedReward() float32 {
	return float32(a.Reward()) / float32(a.Trial())
}
