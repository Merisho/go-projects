package test

import (
    "github.com/merisho/tcp-fs-chat/client"
    "github.com/merisho/tcp-fs-chat/server"
    "github.com/stretchr/testify/suite"
    "log"
    "sync"
    "testing"
    "time"
)

const (
    host = "localhost"
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
    testMessages := []string{ "Message 1", "Message 2", "Message 3", "Message 4", "Message 5" }

    wg := sync.WaitGroup{}

    c1, err := client.ConnectTCP(host, port)
    e2e.NoError(err)
    e2e.NotEmpty(c1.ID())

    c2, err := client.ConnectTCP(host, port)
    e2e.NoError(err)
    e2e.NotEmpty(c2.ID())

    c3, err := client.ConnectTCP(host, port)
    e2e.NoError(err)
    e2e.NotEmpty(c2.ID())

    clients := []client.Client{c1, c2, c3}
    numberOfOtherParticipants := len(clients) - 1
    expectedNumberOfMessagesForEachClient := (len(clients) - 1) * len(testMessages)

    wg.Add(len(clients))

    for _, c := range clients {
        go func(c client.Client) {
            for _, msg := range testMessages {
                err := c.Send(msg)
                e2e.NoError(err)
            }

            msgs := make(map[string]int)
            for i := 0; i < expectedNumberOfMessagesForEachClient; i++ {
                msgs[<- c.Receive()]++
            }

            for _, msg := range testMessages {
                e2e.Equal(numberOfOtherParticipants, msgs[msg])
            }

            wg.Done()
        }(c)
    }

    wg.Wait()
}

func (e2e *E2ETestSuite) TestVeryLongMessageSplitIntoMultipleMessages() {
    testMessage := make([]byte, 2 * 8)
    for i := range testMessage {
        testMessage[i] = 'a' + byte(i % 27)
    }

    c1, err := client.ConnectTCP(host, port)
    e2e.NoError(err)

    c2, err := client.ConnectTCP(host, port)
    e2e.NoError(err)

    err = c1.Send(string(testMessage))
    e2e.NoError(err)

    r := c2.Receive()

    var actualMessages []byte
    for i := 1; i <= 2; i++ {
        select {
        case msg := <- r:
            actualMessages = append(actualMessages, msg...)
        case <- time.After(100 * time.Millisecond):
            e2e.FailNowf("message has not been received", "message %d", i)
        }
    }

    e2e.Equal(testMessage, actualMessages)
}
