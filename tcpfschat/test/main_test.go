package test

import (
    "github.com/merisho/tcp-fs-chat/client"
    "github.com/merisho/tcp-fs-chat/server"
    "github.com/stretchr/testify/suite"
    "log"
    "sync"
    "testing"
)

const (
    port = 1337
)

func TestE2E(t *testing.T) {
    suite.Run(t, new(E2ETestSuite))
}

type E2ETestSuite struct {
    suite.Suite
    server *server.Server
}

func (e2e *E2ETestSuite) SetupSuite() {
    srv, err := server.ServeTCP(port)
    if err != nil {
        log.Fatal(err)
    }

    e2e.server = srv
}

func (e2e *E2ETestSuite) TestE2E() {
    testMessages := []string{ "1", "2", "3", "4", "5" }

    wg := sync.WaitGroup{}

    c1, err := client.ConnectTCP("localhost", port)
    e2e.NoError(err)
    e2e.NotEmpty(c1.ID())

    c2, err := client.ConnectTCP("localhost", port)
    e2e.NoError(err)
    e2e.NotEmpty(c2.ID())

    clients := []client.Client{c1, c2}

    wg.Add(len(clients))

    for _, c := range clients {
        go func(c client.Client) {
            for _, msg := range testMessages {
                err := c.Send(msg)
                e2e.NoError(err)
            }

            var msgs []string
            for range testMessages {
                msgs = append(msgs, <- c.Receive())
            }

            for _, msg := range testMessages {
                e2e.Contains(msgs, msg)
            }

            wg.Done()
        }(c)
    }

    wg.Wait()
}

func (e2e *E2ETestSuite) TestDisconnectClient() {
   c, err := client.ConnectTCP("localhost", port)
   e2e.NoError(err)

   err = e2e.server.Disconnect(c.ID())
   e2e.NoError(err)

   e2e.Equal(0, e2e.server.ConnectionCount())

   _, ok := <- c.Receive()
   e2e.False(ok)
}
