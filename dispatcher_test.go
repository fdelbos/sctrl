package sctrl_test

import (
	"bytes"
	"context"
	"time"

	"github.com/dchest/uniuri"
	. "github.com/fdelbos/sctrl"
	"github.com/fdelbos/sctrl/mocks"
	"github.com/stretchr/testify/mock"
)

func (s *Suite) TestDispatcher() {
	ctx, cancel := context.WithCancel(s.ctx)
	input := mocks.NewDispatchReader(s.T())
	handler := mocks.NewNotificationHandler(s.T())
	var dispatcher Dispatcher
	chRead := make(chan string, 1)

	input.
		On("Read", mock.Anything).
		Return(func(dest []byte) (int, error) {
			msg := <-chRead
			buff := bytes.NewBuffer(dest)
			buff.Reset()
			buff.WriteString(msg)
			return len(buff.Bytes()), nil
		})

	var onMessage func(msg string)
	handler.
		On("OnMessage", mock.Anything).
		Run(func(args mock.Arguments) {
			onMessage(args.String(0))
		})

	s.Run("should dispatch notification on incoming", func() {
		gotNofitication := make(chan struct{}, 1)
		defer close(gotNofitication)
		msg1 := uniuri.New()

		onMessage = func(msg string) {
			s.Require().Equal(msg, msg1)
			gotNofitication <- struct{}{}
		}

		dispatcher = NewDispatcher(ctx, input, handler)
		s.Require().NotNil(dispatcher)
		chRead <- msg1
		<-gotNofitication
	})

	s.Run("should dispatch response on request", func() {
		gotResponse := make(chan struct{}, 1)
		defer close(gotResponse)
		msg := uniuri.New()

		onMessage = func(msg string) {
			s.Fail("should not receive notification")
		}

		go func() {
			dispatcher.Request(time.Second, func(rr ResponseReader) error {
				chRead <- msg
				res, err := rr.GetLine()
				s.Require().NoError(err)
				s.Require().Equal(msg, res)
				gotResponse <- struct{}{}
				return nil
			})
		}()

		<-gotResponse
	})

	s.Run("request multiline", func() {
		gotResponse := make(chan struct{}, 1)
		defer close(gotResponse)
		msg1 := uniuri.New()
		msg2 := uniuri.New()
		msg3 := uniuri.New()

		onMessage = func(msg string) {
			s.Fail("should not receive notification")
		}

		go func() {
			dispatcher.Request(time.Millisecond*100, func(rr ResponseReader) error {
				chRead <- msg1 + "\r\n" + msg2 + "\r\n"
				line1, err := rr.GetLine()
				s.Require().NoError(err)
				s.Require().Equal(msg1, line1)

				line2, err := rr.GetLine()
				s.Require().NoError(err)
				s.Require().Equal(msg2, line2)

				chRead <- msg3 + "\r\n"
				line3, err := rr.GetLine()
				s.Require().NoError(err)
				s.Require().Equal(msg3, line3)

				gotResponse <- struct{}{}
				return nil
			})
		}()

		<-gotResponse
	})

	s.Run("request assert", func() {
		gotResponse := make(chan struct{}, 1)
		defer close(gotResponse)
		msg1 := uniuri.New()
		msg2 := uniuri.New()

		onMessage = func(msg string) {
			s.Fail("should not receive notification")
		}

		go func() {
			dispatcher.Request(
				time.Millisecond*100,
				func(rr ResponseReader) error {
					chRead <- msg1 + "\r\n" + msg2 + "\r\n"
					err := rr.AssertLine(msg1)
					s.Require().NoError(err)

					err = rr.AssertLine("invalid")
					var unexpected *ErrUnexpectedResponse
					s.Require().ErrorAs(err, &unexpected)
					s.Require().Equal(msg2, unexpected.Received)
					s.Require().Equal("invalid", unexpected.Expected)

					gotResponse <- struct{}{}
					return nil
				})
		}()

		<-gotResponse
	})

	s.Run("request recycle", func() {
		gotResponse := make(chan struct{}, 1)
		defer close(gotResponse)
		msg1 := uniuri.New()
		msg2 := uniuri.New()
		msg3 := uniuri.New()

		count := 0
		onMessage = func(msg string) {
			if count == 0 {
				s.Require().Equal(msg2, msg)
				count++
			} else if count == 1 {
				s.Require().Equal(msg3, msg)
				count++
				gotResponse <- struct{}{}
			} else {
				s.Fail("should not receive notification")
			}
		}

		go func() {
			dispatcher.Request(time.Millisecond*100, func(rr ResponseReader) error {
				chRead <- msg1 + "\r\n" + msg2 + "\r\n" + msg3 + "\r\n"
				line1, err := rr.GetLine()
				s.Require().NoError(err)
				s.Require().Equal(msg1, line1)
				return nil
			})
		}()

		<-gotResponse
	})

	s.Run("request timeout", func() {
		gotResponse := make(chan struct{}, 1)
		defer close(gotResponse)

		onMessage = func(msg string) {
			s.Fail("should not receive notification")
		}

		start := time.Now()
		go func() {
			dispatcher.Request(time.Millisecond*100, func(rr ResponseReader) error {

				_, err := rr.GetLine()
				s.Require().ErrorIs(err, ErrTimeout)
				gotResponse <- struct{}{}
				return nil
			})
		}()

		<-gotResponse
		s.Require().True(time.Since(start) > time.Millisecond*100)
		s.Require().True(time.Since(start) < time.Millisecond*150)
	})

	s.Run("request discards empty lines", func() {
		gotResponse := make(chan struct{}, 1)
		defer close(gotResponse)

		onMessage = func(msg string) {
			s.Fail("should not receive notification")
		}

		go func() {
			dispatcher.Request(time.Millisecond*100, func(rr ResponseReader) error {
				chRead <- "\r\n\r\n\r\n"

				// if we dont send the newline, the dispatcher will block
				_, err := rr.GetLine()
				s.Require().ErrorIs(err, ErrTimeout)
				gotResponse <- struct{}{}
				return nil
			})
		}()

		<-gotResponse
	})

	s.Run("close", func() {
		cancel()
	})
}
