package paxosgo

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Learner struct {
	//net
	lis net.Listener
	//learner id
	id int
	//record acceptor had accepted proposal : [acceptor's id ] message
	acceptedMsg map[int]MsgArgs
}

func newLearner(id int, allAcceptor []int) *Learner {
	learner := &Learner{
		id:          id,
		acceptedMsg: make(map[int]MsgArgs),
	}
	for _, aid := range allAcceptor {
		learner.acceptedMsg[aid] = MsgArgs{
			ProposalNumber: 0,
			ProposalValue:  nil,
		}
	}
	learner.server()
	return learner
}

func (l *Learner) Learn(args *MsgArgs, reply *MsgReply) error {
	msg := l.acceptedMsg[args.From]
	if msg.ProposalNumber < args.ProposalNumber {
		l.acceptedMsg[args.From] = *args
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

func (l *Learner) chosen() interface{} {
	acceptCount := make(map[int]int)
	acceptMsg := make(map[int]MsgArgs)
	for _, acceptd := range l.acceptedMsg {
		if acceptd.ProposalNumber != 0 {
			acceptCount[acceptd.ProposalNumber]++
			acceptMsg[acceptd.ProposalNumber] = acceptd
		}
	}

	for n, count := range acceptCount {
		if count >= l.majority() {
			return acceptMsg[n].ProposalValue
		}
	}
	return nil
}

func (l *Learner) server() {
	server := rpc.NewServer()
	server.Register(l)
	addr := fmt.Sprintf(":%d", l.id)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	l.lis = listener
	go func() {
		for {
			conn, err := l.lis.Accept()
			if err != nil {
				continue
			}
			go server.ServeConn(conn)
		}
	}()
}

func (l *Learner) close() {
	l.lis.Close()
}

func (l *Learner) majority() int {
	return len(l.acceptedMsg)/2 + 1
}
