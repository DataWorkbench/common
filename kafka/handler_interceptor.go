package kafka

import "context"

type handlerInterceptor func(ctx context.Context, messages []*ConsumerMessage, handler MessageHandler) (err error)

func chainInterceptors(interceptors []handlerInterceptor) handlerInterceptor {
	var interceptor handlerInterceptor

	if len(interceptors) == 0 {
		interceptor = nil
	} else if len(interceptors) == 1 {
		interceptor = interceptors[0]
	} else {
		interceptor = func(ctx context.Context, messages []*ConsumerMessage, handler MessageHandler) (err error) {
			return interceptors[0](ctx, messages, getChainHandler(interceptors, 0, handler))
		}
	}
	return interceptor
}

func getChainHandler(interceptors []handlerInterceptor, curr int, finalHandler MessageHandler) MessageHandler {
	if curr == len(interceptors)-1 {
		return finalHandler
	}
	return func(ctx context.Context, messages []*ConsumerMessage) (err error) {
		return interceptors[curr+1](ctx, messages, getChainHandler(interceptors, curr+1, finalHandler))
	}
}
