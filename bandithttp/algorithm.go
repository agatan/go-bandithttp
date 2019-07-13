package bandithttp

import (
	"math/rand"
	"time"
)

type Algorithm interface {
	Choice(arms []*Arm) *Arm
}

type AlgorithmFunc func(arms []*Arm) *Arm

func (f AlgorithmFunc) Choice(arms []*Arm) *Arm {
	return f(arms)
}

type epsilonGreedy struct {
	random  *rand.Rand
	epsilon float32
}

func (e *epsilonGreedy) Choice(arms []*Arm) *Arm {
	if e.random.Float32() < e.epsilon {
		// Explore
		n := e.random.Intn(len(arms))
		return arms[n]
	}
	// Exploit
	var (
		maxArm            *Arm
		maxExpectedReward float32 = -1.0
	)
	for _, arm := range arms {
		er := arm.ExpectedReward()
		if maxArm == nil || er > maxExpectedReward || (er == maxExpectedReward && e.random.Float32() < 0.5) {
			maxArm = arm
			maxExpectedReward = er
		}
	}
	return maxArm
}

func EpsilonGreedy(eps float32) Algorithm {
	return &epsilonGreedy{
		random:  rand.New(rand.NewSource(time.Now().UnixNano())),
		epsilon: eps,
	}
}
