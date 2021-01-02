//Package id contains the set of middleware functions
package mid

import  (
    "log"
    "time"
    "context"
    "net/http"
    "github.com/yaowenqiang/service/foundation/web"
)


func Logger(log *log.Logger) web.Middleware {
    m := func(beforeAfter web.Handler) web.Handler {

       h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
           v, ok := ctx.Value(web.KeyValues).(*web.Values)
           if !ok {
                return web.NewShutdownError("web value missing from context")
           }
           log.Printf("%s : started:   %s %s -> %s",
               v.TraceID,
               r.Method, r.URL.Path, r.RemoteAddr,
           )
           err := beforeAfter(ctx, w, r)
           log.Printf("%s : completed: %s %s -> %s (%d) (%s)",
               v.TraceID,
               r.Method, r.URL.Path, r.RemoteAddr,
               v.StatusCode, time.Since(v.Now),
           )
            return err
       }
        return h
    }
    return m

}
