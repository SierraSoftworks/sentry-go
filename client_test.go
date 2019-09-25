package sentry

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleClient() {
	// You can create a new root client directly and configure
	// it by passing any options you wish
	cl := NewClient(
		DSN(""),
	)

	// You can then create a derivative client with any context-specific
	// options. These are useful if you want to encapsulate context-specific
	// information like the HTTP request that is being handled.
	var r *http.Request
	ctxCl := cl.With(
		HTTPRequest(r).WithHeaders(),
		Logger("http"),
	)

	// You can then use the client to capture an event and send it to Sentry
	err := fmt.Errorf("an error occurred")
	ctxCl.Capture(
		ExceptionForError(err),
	)
}

func ExampleDefaultClient() {
	DefaultClient().Capture(
		Message("This is an example message"),
	)
}

func TestDefaultClient(t *testing.T) {
	cl := DefaultClient()
	assert.NotNil(t, cl, "it should not return nil")
	assert.Implements(t, (*Client)(nil), cl, "it should implement the Client interface")
	assert.Equal(t, defaultClient, cl, "it should return the global defaultClient")
}

func TestNewClient(t *testing.T) {
	opt := &testOption{}
	cl := NewClient(opt)

	if assert.NotNil(t, cl, "it should not return nil") {
		assert.Implements(t, (*Client)(nil), cl, "it should implement the Client interface")

		cll, ok := cl.(*client)
		assert.True(t, ok, "it should actually return a *client")

		assert.Nil(t, cll.parent, "it should set the parent of the client to nil")
		assert.Equal(t, []Option{opt}, cll.options, "it should set the client's options correctly")

		t.Run("Capture()", func(t *testing.T) {
			tr := testNewTestTransport()
			assert.NotNil(t, tr, "the test transport should not be nil")

			cl := NewClient(UseTransport(tr))
			assert.NotNil(t, cl, "the client should not be nil")

			e := cl.Capture(Message("test"))
			assert.NotNil(t, e, "the event handle should not be nil")

			ei, ok := e.(QueuedEventInternal)
			assert.True(t, ok, "the event handle should be convertible to an internal queued event")

			select {
			case p := <-tr.ch:
				assert.Equal(t, ei.Packet(), p, "the packet should match the internal event's packet")
				
				pi, ok := p.(*packet)
				assert.True(t, ok, "the packet should actually be a *packet")
				assert.Contains(t, *pi, Message("test").Class(), "the packet should contain the message")
				assert.Equal(t, Message("test"), (*pi)[Message("test").Class()], "the message should be serialized under its key")
			case <-time.After(100 * time.Millisecond):
				t.Error("the event was not dispatched within the timeout of 100ms")
			}
		})

		t.Run("With()", func(t *testing.T) {
			opt := &testOption{}

			ctxCl := cl.With(opt)
			assert.NotNil(t, ctxCl, "the new client should not be nil")
			assert.NotEqual(t, cl, ctxCl, "the new client should not be the same as the old client")

			cll, ok := ctxCl.(*client)
			assert.True(t, ok, "the new client should actually be a *client")
			assert.Equal(t, cl, cll.parent, "the new client should have its parent configured to be the original client")
			assert.Equal(t, []Option{opt}, cll.options, "the new client should have the right list of options")
		})

		t.Run("GetOption()", func(t *testing.T) {
			assert.Nil(t, NewClient().GetOption("unknown-option-class"), "it should return nil for an unrecognized option")
			assert.Nil(t, NewClient(nil).GetOption("unknown-option-class"), "it should ignore nil options")

			opt := &testOption{}
			assert.Equal(t, opt, NewClient(opt).GetOption("test"), "it should return an option if it is present")
			assert.Equal(t, opt, NewClient(&testOption{}, opt).GetOption("test"), "it should return the most recent non-mergeable option")

			assert.Equal(t, &testMergeableOption{3}, NewClient(&testMergeableOption{1}, &testMergeableOption{2}).GetOption("test"), "it should merge options when they support it")
		})

		t.Run("fullDefaultOptions()", func(t *testing.T) {
			opts := cll.fullDefaultOptions()
			assert.NotNil(t, opts, "the full options list should not be nil")
			assert.NotEmpty(t, opts, "the full options list should not be empty")

			assert.Contains(t, opts, opt, "it should include the options passed to the client")

			i := 0
			for _, provider := range defaultOptionProviders {
				opt := provider()
				if opt == nil {
					continue
				}

				if i >= len(opts) {
					t.Error("there are fewer options than there are providers which return option values")
					break
				}

				assert.IsType(t, opt, opts[i], "Expected opts[%d] to have type %s but got %s instead", i, opt.Class(), opts[i].Class())

				i++
			}

			opt1 := &testMergeableOption{data: 1}
			opt2 := &testMergeableOption{data: 2}
			cl := NewClient(opt1)
			assert.NotNil(t, cl, "the client should not be nil")

			dcl := cl.With(opt2)
			assert.NotNil(t, dcl, "the derived client should not be nil")

			cll, ok := dcl.(*client)
			assert.True(t, ok, "the derived client should actually be a *client")

			opts = cll.fullDefaultOptions()
			assert.Contains(t, opts, opt1, "the parent's options should be present in the list")
			assert.Contains(t, opts, opt2, "the derive client's options should be present in the list")
		})
	}
}

func TestClientConfigInterface(t *testing.T) {
	t.Run("DSN()", func (t *testing.T) {
		t.Run("with no default DSN", func(t *testing.T) {
			oldDefaultOptionProviders := defaultOptionProviders
			defer func() {
				defaultOptionProviders = oldDefaultOptionProviders
			}()

			defaultOptionProviders = []func() Option{}

			cl := NewClient()
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.Equal(t, "", cfg.DSN(), "the client should return an empty DSN")
		})

		t.Run("with a custom DSN option implementation", func(t *testing.T) {
			cl := NewClient(&testCustomClassOption{"sentry-go.dsn"})
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.Equal(t, "", cfg.DSN(), "the client should return an empty DSN")
		})

		t.Run("with no custom DSN", func( t *testing.T) {
			cl := NewClient()
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.Equal(t, "", cfg.DSN(), "the client should return an empty DSN")
		})

		t.Run("with a custom DSN", func( t *testing.T) {
			cl := NewClient(DSN("test"))
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.Equal(t, "test", cfg.DSN(), "the client should return the DSN")
		})

		t.Run("with a custom DSN on the parent", func( t *testing.T) {
			cl := NewClient(DSN("test"))
			assert.NotNil(t, cl, "the client should not be nil")

			dcl := cl.With()
			assert.NotNil(t, dcl, "the derived client should not be nil")

			cfg, ok := dcl.(Config)
			assert.True(t, ok, "the derived client should implement the Config interface")
			assert.Equal(t, "test", cfg.DSN(), "the derived client should return the DSN")
		})
	})

	t.Run("SendQueue()", func(t *testing.T) {
		t.Run("with no default send queue", func(t *testing.T) {
			oldDefaultOptionProviders := defaultOptionProviders
			defer func() {
				defaultOptionProviders = oldDefaultOptionProviders
			}()

			defaultOptionProviders = []func() Option{}

			cl := NewClient()
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.SendQueue(), "the client should not have a nil send queue")
			assert.IsType(t, NewSequentialSendQueue(0), cfg.SendQueue(), "the client should default to the sequential send queue")
		})

		t.Run("with a custom send queue option implementation", func(t *testing.T) {
			cl := NewClient(&testCustomClassOption{"sentry-go.sendqueue"})
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.SendQueue(), "the client should not have a nil send queue")
			assert.IsType(t, NewSequentialSendQueue(0), cfg.SendQueue(), "the client should default to the sequential send queue")
		})

		t.Run("with no custom send queue", func( t *testing.T) {
			cl := NewClient()
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.SendQueue(), "the client should not have a nil send queue")
			assert.IsType(t, NewSequentialSendQueue(0), cfg.SendQueue(), "the client should default to the sequential send queue")
			assert.Equal(t, DefaultClient().GetOption("sentry-go.sendqueue").(*sendQueueOption).queue, cfg.SendQueue(), "the client should use the global default send queue")
		})

		t.Run("with a custom send queue", func( t *testing.T) {
			q := NewSequentialSendQueue(0)
			cl := NewClient(UseSendQueue(q))
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.SendQueue(), "the client should not have a nil send queue")
			assert.Equal(t, q, cfg.SendQueue(), "the client should use the configured send queue")
		})

		t.Run("with a custom send queue on the parent", func( t *testing.T) {
			q := NewSequentialSendQueue(0)
			cl := NewClient(UseSendQueue(q))
			assert.NotNil(t, cl, "the client should not be nil")

			dcl := cl.With()
			assert.NotNil(t, dcl, "the derived client should not be nil")

			cfg, ok := dcl.(Config)
			assert.True(t, ok, "the derived client should implement the Config interface")
			assert.NotNil(t, cfg.SendQueue(), "the client should not have a nil send queue")
			assert.Equal(t, q, cfg.SendQueue(), "the client should use the configured send queue")
		})
	})

	t.Run("Transport()", func(t *testing.T) {
		t.Run("with no default transports", func(t *testing.T) {
			oldDefaultOptionProviders := defaultOptionProviders
			defer func() {
				defaultOptionProviders = oldDefaultOptionProviders
			}()

			defaultOptionProviders = []func() Option{}

			cl := NewClient()
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.Transport(), "the client should not have a nil transport")
			assert.IsType(t, newHTTPTransport(), cfg.Transport(), "the client should default to the HTTP transport")
		})

		t.Run("with a custom transport option implementation", func(t *testing.T) {
			cl := NewClient(&testCustomClassOption{"sentry-go.transport"})
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.Transport(), "the client should not have a nil transport")
			assert.IsType(t, newHTTPTransport(), cfg.Transport(), "the client should default to the HTTP transport")
		})

		t.Run("with no custom transport", func( t *testing.T) {
			cl := NewClient()
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.Transport(), "the client should not have a nil transport")
			assert.IsType(t, newHTTPTransport(), cfg.Transport(), "the client should default to the HTTP transport")
			assert.Equal(t, DefaultClient().GetOption("sentry-go.transport").(*transportOption).transport, cfg.Transport(), "the client should use the global default transport")
		})

		t.Run("with a custom transport", func( t *testing.T) {
			tr := newHTTPTransport()
			cl := NewClient(UseTransport(tr))
			assert.NotNil(t, cl, "the client should not be nil")

			cfg, ok := cl.(Config)
			assert.True(t, ok, "the client should implement the Config interface")
			assert.NotNil(t, cfg.Transport(), "the client should not have a nil transport")
			assert.Equal(t, tr, cfg.Transport(), "the client should use the configured transport")
		})

		t.Run("with a custom transport on the parent", func( t *testing.T) {
			tr := newHTTPTransport()
			cl := NewClient(UseTransport(tr))
			assert.NotNil(t, cl, "the client should not be nil")

			dcl := cl.With()
			assert.NotNil(t, dcl, "the derived client should not be nil")

			cfg, ok := dcl.(Config)
			assert.True(t, ok, "the derived client should implement the Config interface")
			assert.NotNil(t, cfg.Transport(), "the client should not have a nil transport")
			assert.Equal(t, tr, cfg.Transport(), "the client should use the configured transport")
		})
	})
}
