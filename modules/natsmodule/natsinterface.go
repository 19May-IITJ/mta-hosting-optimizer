package natsmodule

import (
	"time"

	"github.com/nats-io/nats.go"
)

type NATSConnInterface interface {
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	Publish(subj string, data []byte) error
	Close()
	// Add any other methods that RetrieveHostnames uses from nats.Conn
}
