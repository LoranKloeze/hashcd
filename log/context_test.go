package log

import (
	"context"
	"testing"
)

func TestLoggerContext(t *testing.T) {
	ctx := context.Background()

	ctx = WithLogger(ctx, G(ctx).WithField("test", "one"))
	if G(ctx).Data["test"] != "one" {
		t.Errorf("Expected test field to be one, got %s", G(ctx).Data["test"])
	}

	if G(ctx) != GetLogger(ctx) {
		t.Errorf("Expected the same logger, got different ones")
	}
}
