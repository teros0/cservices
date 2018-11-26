package ingestor

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/teros0/cservices/resources"
)

type Ingestor struct {
	client resources.StorageClient
}

func (i *Ingestor) Run(ctx context.Context, path string, c resources.StorageClient) (err error) {
	i.client = c
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("couldn't open file %s -> %s", path, err)
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
	c := context.Background()
	_, err = i.client.SaveRecord(c,
		&resources.Record{
			Id:    id,
			Name:  name,
			Email: email,
			Phone: phone,
		})
	if err != nil {
		return fmt.Errorf("while saving record {%s, %s, %s, %s} to storage -> %s", id, name, email, phone, err)
	}
	return nil
}
