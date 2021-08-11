package workers

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
)

type NATSWorker struct {
	Nats *nats.EncodedConn
}

// Run should run start processing the worker and be a blocking operation.
func (n *NATSWorker) Run(ctx context.Context) error {
	log.Println("starting nats worker")
	_, err := n.Nats.Conn.QueueSubscribe("orders.new", "worker", handleMessage)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

// Halt should tell the worker to stop doing work.
func (n *NATSWorker) Halt(_ context.Context) error {
	log.Println("stopping nats worker")
	if err := n.Nats.Drain(); err != nil {
		return err
	}
	n.Nats.Close()
	return nil
}

func handleMessage(msg *nats.Msg) {
	log.Println("received new order")
	log.Println(string(msg.Data))
}
