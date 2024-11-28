package paxosgo

import "testing"

func start(acceptorId []int, learnerId []int) ([]*Acceptor, []*Learner) {
	acceptors := make([]*Acceptor, 0)
	for _, aid := range acceptorId {
		a := newAcceptor(aid, learnerId)
		acceptors = append(acceptors, a)
	}

	learners := make([]*Learner, 0)
	for _, lid := range learnerId {
		l := newLearner(lid, acceptorId)
		learners = append(learners, l)
	}

	return acceptors, learners
}

func cleanUp(ac []*Acceptor, le []*Learner) {
	for _, a := range ac {
		a.close()
	}
	for _, l := range le {
		l.close()
	}
}

func TestSingleProposer(t *testing.T) {
	acceptorId := []int{1001, 1002, 1003}
	learnerId := []int{2001}
	acceptors, learners := start(acceptorId, learnerId)
	defer cleanUp(acceptors, learners)

	p := &Proposer{
		id:        1,
		acceptors: acceptorId,
	}
	value := p.propose("hello  world")
	if value != "hello  world" {
		t.Errorf("value = %s , expected = %s\n", value, "hello  world")
	}
	learnValue := learners[0].chosen()
	if learnValue != value {
		t.Errorf("learnValue = %s , expected = %s\n", learnValue, value)
	}
}

func TestTwoProposer(t *testing.T) {
	acceptorId := []int{1001, 1002, 1003}
	learnerId := []int{2001}
	acceptors, learners := start(acceptorId, learnerId)
	defer cleanUp(acceptors, learners)

	p1 := &Proposer{
		id:        1,
		acceptors: acceptorId,
	}
	p2 := &Proposer{
		id:        2,
		acceptors: acceptorId,
	}
	value1 := p1.propose("hello  world")
	value2 := p2.propose("hello  paxos")

	if value1 != value2 {
		t.Errorf("value1 = %s , value2= %s\n", value1, value2)
	}
	learnValue := learners[0].chosen()
	if learnValue != value1 {
		t.Errorf("learnValue = %s , expected = %s\n", learnValue, value1)
	}
}
