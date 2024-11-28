package paxosgo

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Acceptor struct {
	//net
	lis net.Listener
	//server id
	id int
	//promise proposal's id (0 mean doesn't receive message)
	promiseId int
	//accept proposal's id (0 mean doesn't accept message)
	acceptId int
	//accept proposal's value
	acceptValue interface{}
	//learner id list
	learners []int
}

func newAcceptor(id int, learners []int) *Acceptor {
	ac := &Acceptor{
		id:       id,
		learners: learners,
	}
	ac.server()
	return ac
}

// first RPC:
func (ac *Acceptor) Prepare(args *MsgArgs, reply *MsgReply) error {
	if args.ProposalNumber > ac.promiseId {
		ac.promiseId = args.ProposalNumber
		reply.ProposalNumber = ac.acceptId
		reply.ProposalValue = ac.acceptValue
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

// second RPC:
func (ac *Acceptor) Accept(args *MsgArgs, reply *MsgReply) error {
	if args.ProposalNumber >= ac.promiseId {
		ac.promiseId = args.ProposalNumber
		ac.acceptId = args.ProposalNumber
		ac.acceptValue = args.ProposalValue
		reply.Ok = true
		//make learners to learn the proposal
		for _, learnerId := range ac.learners {
			//background opration
			go func(learner int) {
				addr := fmt.Sprintf("127.0.0.1:%d", learner)
				args.From = ac.id
				args.To = learner
				resp := new(MsgReply)
				if ok := call(addr, "Learner.Learn", args, resp); !ok {
					return
				}
			}(learnerId)
		}

	} else {
		reply.Ok = false
	}
	return nil
}

func (ac *Acceptor) server() {
	server := rpc.NewServer()
	server.Register(ac)
	addr := fmt.Sprintf(":%d", ac.id)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	ac.lis = listener
	//
	go func() {
		for {
			conn, err := ac.lis.Accept()
			if err != nil {
				continue
			}
			go server.ServeConn(conn)
		}
	}()
}

func (ac *Acceptor) close() {
	ac.lis.Close()
}
