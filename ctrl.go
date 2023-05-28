package sctrl

import (
	"context"
	"fmt"
	"time"
)

type (
	ErrUnexpectedResponse struct {
		Expected string
		Received string
	}

	NotificationHandler interface {
		OnMessage(string)
	}

	ControllerConnection interface {
		Read([]byte) (int, error)
		Write([]byte) (int, error)
	}

	ResponseReader interface {
		// GetLine reads a line from the response, blocking until a line is available
		// or the context is cancelled.
		GetLine() (string, error)

		// AssertLine reads a line from the response, blocking until a line is available
		// or the context is cancelled. It then asserts that the line matches the
		// expected string.
		AssertLine(string) error

		// AssertOK reads a line containing "OK" from the response, or returns an error
		AssertOK() error

		// ReadLineWithOk reads a new line, immediately followed by an "OK" line.
		ReadLineWithOk() (string, error)
	}

	ResponseHandler func(ResponseReader) error

	Controller interface {
		Command(string, ResponseHandler) error
		CommandWithTimeout(string, time.Duration, ResponseHandler) error
	}

	ctrl struct {
		ctx        context.Context
		cancel     context.CancelFunc
		conn       ControllerConnection
		dispatcher Dispatcher
	}
)

const (
	EOL = "\r\n"
)

var (
	DefaultCommandTimeout = 5 * time.Second

	ErrTimeout = context.DeadlineExceeded
)

func NewController(ctx context.Context, conn ControllerConnection, notificationHandler NotificationHandler) Controller {
	// // no timeout we always want to read
	// conn.SetReadTimeout(serial.NoTimeout)

	ctx, cancel := context.WithCancel(ctx)
	m := &ctrl{
		ctx:        ctx,
		cancel:     cancel,
		conn:       conn,
		dispatcher: NewDispatcher(ctx, conn, notificationHandler),
	}
	return m
}

func (c *ctrl) Command(cmd string, handler ResponseHandler) error {
	return c.CommandWithTimeout(cmd, DefaultCommandTimeout, handler)
}

func (c *ctrl) CommandWithTimeout(cmd string, timeout time.Duration, handler ResponseHandler) error {
	err := c.dispatcher.Request(timeout,
		func(reader ResponseReader) error {
			_, err := c.conn.Write([]byte(cmd + EOL))
			if err != nil {
				return err
			}
			return handler(reader)
		})
	return err
}

func (e *ErrUnexpectedResponse) Error() string {
	return fmt.Sprintf("unexpected response: expected=%s Received=%s", e.Expected, e.Received)
}
