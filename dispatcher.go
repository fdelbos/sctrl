package sctrl

import (
	"bufio"
	"bytes"
	"context"
	"time"
)

type (
	dispatcherToken struct{}

	Dispatcher interface {
		Request(timeout time.Duration, handler ResponseHandler) error
	}

	DispatchReader interface {
		Read([]byte) (int, error)
	}

	dispatcher struct {
		ctx    context.Context
		cancel context.CancelFunc

		// reserve access to the incoming data
		chReserve chan dispatcherToken

		notificationHandler NotificationHandler
		input               DispatchReader
		// data from the input
		chInput chan string
		// response to the request
		chResponse chan string
		// recycle chan
		chRecycle chan string
	}
)

const (
	DispatcherInputBufferSize    = 20
	DispatcherReservedBufferSize = 10
)

func NewDispatcher(ctx context.Context, input DispatchReader, notificationHandler NotificationHandler) Dispatcher {
	ctx, cancel := context.WithCancel(ctx)
	d := dispatcher{
		ctx:                 ctx,
		cancel:              cancel,
		chReserve:           make(chan dispatcherToken, 1),
		chInput:             make(chan string, DispatcherInputBufferSize),
		chResponse:          make(chan string, DispatcherReservedBufferSize),
		chRecycle:           make(chan string, DispatcherReservedBufferSize),
		notificationHandler: notificationHandler,
		input:               input,
	}

	// add the initial token for reserving access to the data
	d.chReserve <- dispatcherToken{}

	go d.dispatch()
	go d.read()

	return &d
}

func (d *dispatcher) recycle() {
	for {
		select {
		case <-d.ctx.Done():
			return

		case line := <-d.chResponse:
			d.chRecycle <- line

		default:
			return
		}
	}
}

func (d *dispatcher) Request(timeout time.Duration, handler ResponseHandler) error {
	// get the token to reserve access to the data
	token := <-d.chReserve
	ctx, cancel := context.WithTimeout(d.ctx, timeout)

	defer func() {
		cancel()
		d.recycle()
		// return the token to the pool
		d.chReserve <- token
	}()

	return handler(&responseReader{
		ctx:        ctx,
		chResponse: d.chResponse,
	})
}

func (d *dispatcher) read() {
	for {
		select {
		case <-d.ctx.Done():
			return

		default:
			// see how we can reuse the buffer
			buff := make([]byte, 2048)
			n, err := d.input.Read(buff)
			if err == nil {
				stream := bytes.NewBuffer(buff[:n])
				scanner := bufio.NewScanner(stream)
				for scanner.Scan() {
					line := scanner.Text()
					if line != "" {
						d.chInput <- line
					}
				}
			} else {
				// TODO: handle error
			}
		}
	}
}

func (d *dispatcher) dispatch() {
	for {
		select {

		case <-d.ctx.Done():
			return

		case line := <-d.chRecycle:
			d.notificationHandler.OnMessage(line)

		case line := <-d.chInput:
			select {
			case token := <-d.chReserve:
				// no request is waiting for data
				// return the token to the pool
				d.chReserve <- token
				d.notificationHandler.OnMessage(line)

			default:
				// a request is waiting for data
				d.chResponse <- line
			}
		}
	}
}
