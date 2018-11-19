package storage

import (
	"context"
	"fmt"
	"net"

	"github.com/teros0/cservices/resources"

	"google.golang.org/grpc"
)

type Storage struct {
	db map[string]interface{}
}

func (s *Storage) SaveRecord(ctx context.Context, r *resources.Record) (e *resources.Empty, err error) {
	s.db[r.Id] = r
	return &resources.Empty{}, nil
}

func (s *Storage) Run(ctx context.Context, port string) (err error) {
	s.db = make(map[string]interface{})

	l, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("couldn't start listener at port %s -> %s", port, err)
	}
	server := grpc.NewServer()
	resources.RegisterStorageServer(server, s)
	err = server.Serve(l)
	if err != nil {
		return fmt.Errorf("while serving -> %s", err)
	}
	return nil
}
