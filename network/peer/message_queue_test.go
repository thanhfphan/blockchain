package peer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhfphan/blockchain/message"
)

func newMsgCreator(t *testing.T) message.Creator {
	t.Helper()

	mc, err := message.NewCreator()
	require.NoError(t, err)

	return mc
}

func Test_BlockingQueue(t *testing.T) {
	q := NewBlockingQueue(10, nil)

	mc := newMsgCreator(t)
	msgs := []message.OutboundMessage{}
	numToSend := 10

	for i := 0; i < numToSend; i++ {
		m, err := mc.Pong("")
		if err != nil {
			t.Error(err)
		}
		msgs = append(msgs, m)
	}

	go func() {
		for i := 0; i < numToSend; i++ {
			result := q.Push(context.Background(), msgs[i])
			if !result {
				t.Error(errors.New("push mesaage to queue failed"))
			}
		}
	}()

	for i := 0; i < numToSend; i++ {
		tMsg, success := q.Pop()
		require.True(t, success)
		require.Equal(t, msgs[i], tMsg)
	}

	_, success := q.PopNow()
	require.False(t, success)
}

func Test_BlockingQueue_CtxCancelledAndQueueIsClosed(t *testing.T) {
	q := NewBlockingQueue(10, nil)

	mc := newMsgCreator(t)

	msg, err := mc.Pong("")
	if err != nil {
		t.Error(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan bool)
	go func() {
		ok := q.Push(ctx, msg)
		require.False(t, ok)
		close(done)
	}()
	<-done

	done = make(chan bool)
	go func() {
		ok := q.Push(context.Background(), msg)
		require.False(t, ok) // false because the Queue was closed
		close(done)
	}()
	q.Close()
	<-done

	_, ok := q.Pop()
	require.False(t, ok)

}
