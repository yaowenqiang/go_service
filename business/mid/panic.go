package mid

import  (
    "log"
    "context"
    "runtime/debug"
    "net/http"
    "github.com/pkg/errors"
    "github.com/yaowenqiang/service/foundation/web"
)


func Panic(log *log.Logger) web.Middleware {
    m := func(handler web.Handler) web.Handler {

       h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
           v, ok := ctx.Value(web.KeyValues).(*web.Values)
           if !ok {
                return web.NewShutdownError("web value missing from context")
           }

           defer func() {
               if r := recover(); r != nil {
                   err  = errors.Errorf("panic: %v", r)
                   log.Printf("%s : PANIC:     \n%s", v.TraceID, debug.Stack())
               }
           }()

           return handler(ctx, w, r)
       }
        return h
    }
    return m

}

