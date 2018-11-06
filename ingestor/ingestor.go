package ingestor

import (
	"bufio"
	"context"
	"cservices/resources"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type Ingestor struct {
}

var (
	serviceAddress = "localhost"
)

func (i *Ingestor) Run(ctx context.Context, path string, port string) (err error) {
	serviceAddress += port
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't opend file %s -> %s", path, err)
	}
	r := csv.NewReader(bufio.NewReader(f))

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		rec, err := r.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("while reading from file -> %s", err)
		}
		if err = i.SendRecord(rec); err != nil {
			return fmt.Errorf("while sending rec -> %s", err)
		}
	}
	return nil
}

func (i *Ingestor) SendRecord(rec []string) (err error) {
	id, name, email, phone := rec[0], rec[1], rec[2], rec[3]
	phone = strings.Join([]string{"+44", phone}, "")
	idi, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("csv contains incorrect data %s -> %s", id, err)
	}
	grpcConn, err := grpc.Dial(
		serviceAddress,
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("while dialing storage service -> %s", err)
	}
	defer grpcConn.Close()

	cl := resources.NewStorageClient(grpcConn)
	c := context.Background()
	_, err = cl.SaveRecord(c,
		&resources.Record{
			Id:    idi,
			Name:  name,
			Email: email,
			Phone: phone,
		})
	if err != nil {
		fmt.Errorf("while saving record {%s, %s, %s, %s} to storage -> %s", id, name, email, phone, err)
	}
	return nil
}
