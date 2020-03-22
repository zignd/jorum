package jorum

import (
	"fmt"
	"io"
)

var jorum = make(map[string]interface{})
var servicesErrorChs = make(map[string]chan error)
var servicesWarnChs = make(map[string]chan error)
var servicesCloseChs = make(map[string]chan error)
var servicesInfoChs = make(map[string]chan string)

var errorChs []chan ErrorEvent
var warnChs []chan ErrorEvent
var closeChs []chan ErrorEvent
var infoChs []chan InfoEvent

var abortCh = make(chan struct{})

// Ready is expected to be called when the services are all set inside the jorum.
func Ready() {
	for name, ch := range servicesErrorChs {
		go func(name string, ch chan error) {
			for {
				select {
				case ev := <-ch:
					emitError(ErrorEvent{
						Name:  name,
						Error: ev,
					})
				case <-abortCh:
					return
				}
			}
		}(name, ch)
	}
	for name, ch := range servicesWarnChs {
		go func(name string, ch chan error) {
			for {
				select {
				case ev := <-ch:
					emitWarn(ErrorEvent{
						Name:  name,
						Error: ev,
					})
				case <-abortCh:
					return
				}
			}
		}(name, ch)
	}
	for name, ch := range servicesCloseChs {
		go func(name string, ch chan error) {
			for {
				select {
				case ev := <-ch:
					emitClose(ErrorEvent{
						Name:  name,
						Error: ev,
					})
				case <-abortCh:
					return
				}
			}
		}(name, ch)
	}
	for name, ch := range servicesInfoChs {
		go func(name string, ch chan string) {
			for {
				select {
				case ev := <-ch:
					emitInfo(InfoEvent{
						Name:    name,
						Message: ev,
					})
				case <-abortCh:
					return
				}
			}
		}(name, ch)
	}
}

// Get retrieves a registered service.
func Get(name string) (interface{}, error) {
	value := jorum[name]
	if value == nil {
		return nil, fmt.Errorf("%s not found in the jorum", name)
	}
	return value, nil
}

// GetNoErr retrieves a registered service, returning nil when it's not registered by the provided name.
func GetNoErr(name string) interface{} {
	value, _ := Get(name)
	return value
}

// Register a service by the provided name.
func Register(name string, service interface{}) error {
	_, inUse := jorum[name]
	if inUse {
		return fmt.Errorf("there is already an service named %s in the jorum", name)
	}
	jorum[name] = service
	if on, ok := service.(OnErrorer); ok {
		ch := make(chan error, 100)
		on.OnError(ch)
		servicesErrorChs[name] = ch
	}
	if on, ok := service.(OnWarner); ok {
		ch := make(chan error, 100)
		on.OnWarn(ch)
		servicesWarnChs[name] = ch
	}
	if on, ok := service.(OnCloser); ok {
		ch := make(chan error, 100)
		on.OnClose(ch)
		servicesCloseChs[name] = ch
	}
	if on, ok := service.(OnInfoer); ok {
		ch := make(chan string, 100)
		on.OnInfo(ch)
		servicesInfoChs[name] = ch
	}
	return nil
}

// Close calls on the registered services.
func Close() error {
	close(abortCh)
	for name, service := range jorum {
		c, ok := service.(io.Closer)
		if !ok {
			continue
		}
		emitInfo(InfoEvent{
			Name:    name,
			Message: fmt.Sprintf("jorum is trying to close an service named %s", name),
		})
		if err := c.Close(); err != nil {
			return fmt.Errorf("failed to call Close on jorum service named %s: %w", name, err)
		}
		emitInfo(InfoEvent{
			Name:    name,
			Message: fmt.Sprintf("jorum successfully closed an service named %s", name),
		})
	}
	return nil
}

// OnError receives a channel to be notified whenever an error is emitted by a registered service.
func OnError(ch chan ErrorEvent) {
	errorChs = append(errorChs, ch)
}

func emitError(ev ErrorEvent) {
	for _, ch := range errorChs {
		ch <- ev
	}
}

// OnWarn receives a channel to be notified whenever a warning is emitted by a registered service.
func OnWarn(ch chan ErrorEvent) {
	warnChs = append(warnChs, ch)
}

func emitWarn(ev ErrorEvent) {
	for _, ch := range warnChs {
		ch <- ev
	}
}

// OnClose receives a channel to be notified whenever a registered service closes.
func OnClose(ch chan ErrorEvent) {
	closeChs = append(closeChs, ch)
}

func emitClose(ev ErrorEvent) {
	for _, ch := range closeChs {
		ch <- ev
	}
}

// OnInfo receives a channel to be notified whenever an info event is emitted by a registered service.
func OnInfo(ch chan InfoEvent) {
	infoChs = append(infoChs, ch)
}

func emitInfo(ev InfoEvent) {
	for _, ch := range infoChs {
		ch <- ev
	}
}

// ErrorEvent contains the error and the name of the service that emitted it.
type ErrorEvent struct {
	Name  string
	Error error
}

// InfoEvent contains the message and the name of the service that emitted it.
type InfoEvent struct {
	Name    string
	Message string
}

// OnInfoer is the interface for services intending to emit info events.
type OnInfoer interface {
	OnInfo(chan string)
}

// OnErrorer is the interface for services intending to emit error events.
type OnErrorer interface {
	OnError(chan error)
}

// OnWarner is the interface for services intending to emit warning events.
type OnWarner interface {
	OnWarn(chan error)
}

// OnCloser is the interface for services intending to emit close events.
type OnCloser interface {
	OnClose(chan error)
}
