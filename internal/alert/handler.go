package alert

import "context"

// Handler is implemented by anything that can receive an Alert.
type Handler interface {
	Handle(ctx context.Context, a Alert) error
}

// HandlerFunc is a function adapter for Handler.
type HandlerFunc func(ctx context.Context, a Alert) error

// Handle implements Handler.
func (f HandlerFunc) Handle(ctx context.Context, a Alert) error {
	return f(ctx, a)
}

// LogHandler is a Handler that writes alerts to a standard logger.
type LogHandler struct {
	printf func(format string, args ...any)
}

// NewLogHandler returns a LogHandler that uses the provided printf-style func.
func NewLogHandler(printf func(format string, args ...any)) *LogHandler {
	return &LogHandler{printf: printf}
}

// Handle implements Handler.
func (h *LogHandler) Handle(_ context.Context, a Alert) error {
	h.printf("[portwatch] %s port %d/%s (pid %d)",
		a.Kind, a.Port.Port, a.Port.Proto, a.Port.PID)
	return nil
}
