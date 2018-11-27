package ingestor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teros0/cservices/resources"
	"google.golang.org/grpc"
)

type FakeClient struct{}

func (c *FakeClient) SaveRecord(ctx context.Context, in *resources.Record, opts ...grpc.CallOption) (*resources.Empty, error) {
	return &resources.Empty{}, nil
}

func TestSendRecord(t *testing.T) {
	table := []struct {
		rec []string
		err error
	}{
		{[]string{"1", "2", "3", "4"}, nil},
		{[]string{"1", "2", "3"}, errors.New("test")},
	}
	i := Ingestor{
		client: &FakeClient{},
	}
	for _, v := range table {
		err := i.SendRecord(v.rec)
		assert.IsType(t, v.err, err)
	}
}
