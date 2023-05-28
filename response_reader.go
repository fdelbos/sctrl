package sctrl

import "context"

type (
	responseReader struct {
		ctx        context.Context
		chResponse chan string
	}
)

func (rr *responseReader) GetLine() (string, error) {
	select {
	case <-rr.ctx.Done():
		return "", ErrTimeout
	case line := <-rr.chResponse:
		return line, nil
	}
}

func (rr *responseReader) AssertLine(msg string) error {
	if line, err := rr.GetLine(); err != nil {
		return err
	} else if line != msg {
		return &ErrUnexpectedResponse{
			Expected: msg,
			Received: line,
		}
	}
	return nil
}

func (rr *responseReader) AssertOK() error {
	return rr.AssertLine("OK")
}

func (rr *responseReader) ReadLineWithOk() (string, error) {
	line, err := rr.GetLine()
	if err != nil {
		return "", err
	}
	return line, rr.AssertOK()
}
