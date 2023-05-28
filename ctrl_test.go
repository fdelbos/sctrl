package sctrl_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/fdelbos/sctrl/mocks"

	. "github.com/fdelbos/sctrl"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	ctx context.Context
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}

func (s *Suite) SetupTest() {
	s.ctx = context.Background()
}

func (s *Suite) TestCtrl() {
	conn := mocks.NewControllerConnection(s.T())
	handler := mocks.NewNotificationHandler(s.T())
	ctrl := NewController(s.ctx, conn, handler)
	chRead := make(chan string, 1)
	chWrite := make(chan string, 1)

	conn.
		On("Read", mock.Anything).
		Return(func(dest []byte) (int, error) {
			msg := <-chRead
			buff := bytes.NewBuffer(dest)
			buff.Reset()
			buff.WriteString(msg)
			return len(buff.Bytes()), nil
		})

	conn.
		On("Write", mock.Anything).
		Return(func(data []byte) (int, error) {
			chWrite <- string(data)
			return len(data), nil
		})

	var onMessage func(msg string)
	handler.
		On("OnMessage", mock.Anything).
		Run(func(args mock.Arguments) {
			onMessage(args.String(0))
		})

	s.Run("should send command and read the result", func() {
		cmd := "command"
		resp := "response"
		onMessage = func(msg string) {
			s.Fail("should not receive notification")
		}

		err := ctrl.Command(cmd, func(rr ResponseReader) error {
			wCmd := <-chWrite
			s.Require().Equal(cmd+EOL, wCmd)

			chRead <- resp
			wResp, err := rr.GetLine()
			s.Require().NoError(err)
			s.Require().Equal(resp, wResp)

			return nil
		})
		s.NoError(err)
	})

	s.Run("should timeout on command", func() {
		cmd := "command"
		onMessage = func(msg string) {
			s.Fail("should not receive notification")
		}

		start := time.Now()
		err := ctrl.CommandWithTimeout(
			cmd,
			time.Millisecond*100,
			func(rr ResponseReader) error {
				wCmd := <-chWrite
				s.Require().Equal(cmd+EOL, wCmd)

				_, err := rr.GetLine()
				s.Require().ErrorIs(err, ErrTimeout)
				return nil
			})
		s.NoError(err)
		s.Require().True(time.Since(start) > time.Millisecond*100)
		s.Require().True(time.Since(start) < time.Millisecond*150)
	})

	s.Run("should send notification", func() {
		notif := "notification"
		gotNofitication := make(chan struct{}, 1)
		defer close(gotNofitication)

		onMessage = func(msg string) {
			s.Require().Equal(notif, msg)
			gotNofitication <- struct{}{}
		}

		chRead <- notif
		<-gotNofitication
	})

}
