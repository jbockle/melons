package eventsourcing

import "context"

type contextMetadata struct{}

var ContextMetadataKey = &contextMetadata{}

func WithEventMetadata(ctx context.Context, metadata map[string]any) context.Context {
	if ctx == nil {
		panic("ctx is nil")
	}

	value := ctx.Value(ContextMetadataKey)
	if value != nil {
		for k, v := range metadata {
			value.(map[string]any)[k] = v
		}
	} else {
		value = metadata
	}

	return context.WithValue(ctx, ContextMetadataKey, value)
}

func GetEventMetadata(ctx context.Context) map[string]any {
	if ctx == nil {
		panic("ctx is nil")
	}

	value := ctx.Value(ContextMetadataKey)
	if value == nil {
		return make(map[string]any)
	}

	return value.(map[string]any)
}
