package paxosgo

import (
	"net/rpc"
)

type MsgArgs struct {
	// proposal number
	ProposalNumber int
	// proposal value
	ProposalValue interface{}
	//from server id
	From int
	//to server id
	To int
}

type MsgReply struct {
	//result
	Ok bool
	// proposal number
	ProposalNumber int
	// proposal value
	ProposalValue interface{}
}

func call(srv string, name string, args interface{}, reply interface{}) bool {
	client, err := rpc.Dial("tcp", srv)
	if err != nil {
		return false
	}
	defer client.Close()
	err = client.Call(name, args, reply)
	if err != nil {
		return false
	}
	return true
}
