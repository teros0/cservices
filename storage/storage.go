package storage

import (
	"context"
	"cservices/resources"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Storage struct {
	db map[int64]interface{}
}

func (s *Storage) SaveRecord(ctx context.Context, r *resources.Record) (e *resources.Empty, err error) {
	s.db[r.Id] = r
	return &resources.Empty{}, nil
}

func (s *Storage) Run(ctx context.Context, port string) (err error) {
	s.db = make(map[int64]interface{})

	l, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("couldn't start listener at port %s -> %s", port, err)
	}

	server := grpc.NewServer()
	resources.RegisterStorageServer(server, s)
	if err = server.Serve(l); err != nil {
		return fmt.Errorf("while serving -> %s", err)
	}
	return nil
}
