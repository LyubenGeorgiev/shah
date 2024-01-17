package util

import "context"

func IsAuthenticatedUser(ctx context.Context) bool {
	v, ok := ctx.Value("authenticated").(bool)
	return ok && v
}
