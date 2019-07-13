package bandithttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Experiment struct {
	Name      string
	Arms      []*Arm
	Algorithm Algorithm
}

func (e *Experiment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		e.pull(w, r)
	} else if r.Method == http.MethodPost {
		e.reward(w, r)
	}
}

func (e *Experiment) pull(w http.ResponseWriter, r *http.Request) {
	arm := e.choice()
	arm.trial += 1
	arm.ServeHTTP(w, r)
}

func (e *Experiment) choice() *Arm {
	return e.Algorithm.Choice(e.Arms)
}

func (e *Experiment) GiveReward(name string) {
	for _, arm := range e.Arms {
		if arm.Name == name {
			arm.reward += 1
		}
	}
}

func (e *Experiment) reward(w http.ResponseWriter, r *http.Request) {
	arm, err := e.getArmNameFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	e.GiveReward(arm)
	w.WriteHeader(http.StatusOK)
}

func (e *Experiment) getArmNameFromRequest(r *http.Request) (string, error) {
	var body struct {
		Arm string `json:"arm"`
	}
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(reqBytes, &body); err != nil {
		return "", err
	}
	return body.Arm, nil
}
