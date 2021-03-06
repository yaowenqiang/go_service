package handlers
import (
    "log"
    _ "math/rand"
    "context"
    "os"
    "net/http"
    _ "github.com/pkg/errors"
    "github.com/yaowenqiang/service/foundation/web"
    "github.com/yaowenqiang/service/foundation/database"
    "github.com/jmoiron/sqlx"
)


type checkGroup struct {
    build string
	db *sqlx.DB
    Log *log.Logger
}


//func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error  {
func (cg checkGroup) readiness(ctx context.Context,  w http.ResponseWriter, r *http.Request)  error {
    //if m := rand.Intn(100); m % 2 == 0 {
        //return errors.New("untrusted error")
    //} else {
        //return web.NewRequestError(errors.New("tusted error"), http.StatusNotFound)
        //return web.NewRequestError(errors.New("tusted error"), http.StatusNotFound)
        //panic("forcing panic")
        //return web.NewShutdownError("forcing shutdown")

    //}
        /*status := struct {
            Status string
        }{
            Status: "OK",
        }*/

		status := "OK"
		statusCode := http.StatusOK

		if err := database.StatusCheck(ctx, cg.db); err != nil {
			status = "db not ready"
			statusCode = http.StatusInternalServerError
		}

		health := struct {
			status string  `json:"status"`
		}{
			status: status,
		}

        return web.Respond(ctx, w, health, statusCode)
}


// liveness returns simple status info if the service is alive. if the
// app is deployed to a kubernetes cluster, it will also return pod, node, and
// namespace details via the Downard API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (cg checkGroup) liveness(ctx context.Context,  w http.ResponseWriter, r *http.Request)  error {

    host, err := os.Hostname()
    if err != nil {
        host = "unknown host"
    }

	info := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     cg.build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

    return web.Respond(ctx, w, info, http.StatusOK)
}
