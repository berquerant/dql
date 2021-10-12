package async

import "context"

// IsDone returns true if context is already canceled.
func IsDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
