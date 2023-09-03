package mocking

import (
	"errors"
	"mta2/modules/utility"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

type MockNATSConn struct {
	mock.Mock
}

func (m *MockNATSConn) Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error) {
	args := m.Called(subj, data, timeout)
	return args.Get(0).(*nats.Msg), args.Error(1)
}

func (m *MockNATSConn) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	args := m.Called(subj, cb)
	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (m *MockNATSConn) Publish(subj string, data []byte) error {
	args := m.Called(subj, data)
	return args.Error(0)
}

func (m *MockNATSConn) Close() {
	m.Called()
}

type MockHostingServiceHostMap struct {
	mock.Mock
	data map[string]int // Simulate the data map for testing
}

func NewMockHostMap() *MockHostingServiceHostMap {
	return &MockHostingServiceHostMap{
		data: make(map[string]int),
	}
}

func (m *MockHostingServiceHostMap) IsEmpty() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockHostingServiceHostMap) Set(key string, value string) {
	m.Called(key, value)
}

// Implement the additional functions

func (m *MockHostingServiceHostMap) Put(messages ...*utility.Message) {
	for _, msg := range messages {
		m.data[msg.Hostname] = msg.Active
	}
}

func (m *MockHostingServiceHostMap) Contains(key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *MockHostingServiceHostMap) GetValue(key string) (int, error) {
	value, ok := m.data[key]
	if !ok {
		return 0, errors.New("key not found")
	}
	return value, nil
}

func (m *MockHostingServiceHostMap) GetValues() map[string]int {
	return m.data
}

func (m *MockHostingServiceHostMap) RemoveKey(keys ...string) {
	for _, key := range keys {
		delete(m.data, key)
	}
}

func (m *MockHostingServiceHostMap) Size() int {
	return len(m.data)
}

func (m *MockHostingServiceHostMap) Clear() {
	m.data = make(map[string]int)
}
