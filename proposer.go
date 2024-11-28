package paxosgo

import (
	"fmt"
)

type Proposer struct {
	// server id
	id int
	// max round id
	round int
	// proposal number = (round,id)
	number int
	//acceptor id  list
	acceptors []int
}

func (pro *Proposer) propose(v interface{}) interface{} {
	pro.round++
	pro.number = pro.proposalNumber()

	//phase 1
	prepareNumber := 0
	maxNumber := 0
	for _, acceptor := range pro.acceptors {
		args := &MsgArgs{
			ProposalNumber: pro.number,
			From:           pro.id,
			To:             acceptor,
		}
		reply := new(MsgReply)
		if ok := call(fmt.Sprintf("127.0.0.1:%d", acceptor), "Acceptor.Prepare", args, reply); !ok {
			continue
		}
		if reply.Ok {
			prepareNumber++
			if reply.ProposalNumber > maxNumber {
				maxNumber = reply.ProposalNumber
				v = reply.ProposalValue
			}
		}
		if prepareNumber == pro.majority() {
			break
		}
	}

	//phase 2
	acceptCount := 0
	if prepareNumber >= pro.majority() {
		for _, acceptor := range pro.acceptors {
			args := &MsgArgs{
				ProposalNumber: pro.number,
				ProposalValue:  v,
				From:           pro.id,
				To:             acceptor,
			}
			reply := new(MsgReply)
			if ok := call(fmt.Sprintf("127.0.0.1:%d", acceptor), "Acceptor.Accept", args, reply); !ok {
				continue
			}
			if reply.Ok {
				acceptCount++
			}
		}
	}
	if acceptCount >= pro.majority() {
		return v
	}
	return nil

}

func (pro *Proposer) majority() int {
	return len(pro.acceptors)/2 + 1
}

func (pro *Proposer) proposalNumber() int {
	return pro.round<<16 | pro.id
}
