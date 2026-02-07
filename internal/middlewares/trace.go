package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spinvettle/OctoStudio/internal/consts"
)

func TraceID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var traceID string
		if _, ok := ctx.Get(consts.CtxKeyTraceID); !ok {
			traceID = uuid.NewString()
			ctx.Set(consts.CtxKeyTraceID, traceID)
		}
		//set的key是只在gin内部使用，在orm等服务中没法使用，*gin.Context中Request，
		// 是go的原生request，请求到来的时候go会为这个请求创建request，并放入一个context，
		// 现在我们在创建一个新的contextWithValue放入TraceID，给这个请求链的所有服务使用，
		// 使用Request.WithContext(ctx)替换原来的
		context := context.WithValue(ctx.Request.Context(), consts.CtxKeyTraceID, traceID)
		ctx.Request = ctx.Request.WithContext(context)
		ctx.Next()
	}
}
