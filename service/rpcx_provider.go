package service

import (
	"context"
	"shop/logger"
	"shop/rpc"
)



type Arith int

func (t *Arith) Mul(ctx context.Context, args *rpc.Args, reply *rpc.Reply) error {
	reply.C = args.A * args.B
	logger.Info.Println("rpcx call: %d * %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func (t *Arith) Add(ctx context.Context, args *rpc.Args, reply *rpc.Reply) error {
	reply.C = args.A + args.B
	logger.Info.Println("rpcx call: %d + %d = %d\n", args.A, args.B, reply.C)
	return nil
}

func (t *Arith) Say(ctx context.Context, args *string, reply *string) error {
	*reply = "hello " + *args
	return nil
}

