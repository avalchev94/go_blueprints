package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/avalchev94/go_blueprints/vault"

	"github.com/avalchev94/go_blueprints/vault/client/grpc"

	"google.golang.org/grpc"
)

func main() {
	var grpcAddr = flag.String("addr", ":8081", "grpc server address")
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalln("grpc dial:", err)
	}
	defer conn.Close()
	vaultService := grpcclient.New(conn)

	args := flag.Args()
	var cmd string
	cmd, args = pop(args)
	switch cmd {
	case "hash":
		var password string
		password, args = pop(args)
		hash(ctx, vaultService, password)
	case "validate":
		var password, hash string
		password, args = pop(args)
		hash, args = pop(args)
		validate(ctx, vaultService, password, hash)
	default:
		log.Fatalln("unknown command", cmd)
	}
}

func pop(s []string) (string, []string) {
	if len(s) == 0 {
		return "", s
	}
	return s[0], s[1:]
}

func hash(ctx context.Context, srv vault.Service, password string) {
	h, err := srv.Hash(ctx, password)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(h)
}

func validate(ctx context.Context, srv vault.Service, password, hash string) {
	valid, err := srv.Validate(ctx, password, hash)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(valid)
}
