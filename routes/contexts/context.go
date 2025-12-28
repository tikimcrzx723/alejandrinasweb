package contexts

import (
	"context"
	"fmt"
)

func ExtractApp(ctx context.Context) App {
	appCtx, ok := ctx.Value(AppKey{}).(App)
	if !ok {
		return App{}
	}

	return appCtx
}

func ExtractToken(ctx context.Context) string {
	return ExtractApp(ctx).Token
}

func ExtractRole(ctx context.Context) string {
	return ExtractApp(ctx).Role
}

func ExtractFlashMessages(ctx context.Context) []FlashMessage {
	flashCtx, ok := ctx.Value(FlashKey{}).([]FlashMessage)
	if !ok {
		return []FlashMessage{}
	}

	fmt.Println(flashCtx)

	return flashCtx
}
