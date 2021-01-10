package database

import (
    "net/url"
    "context"
    "strings"
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

type Config struct {
    User string
    Password string
    Host string
    Name string
    DisableTLS bool
}

func Open(cfg Config) (*sqlx.DB, error) {
    sslMode := "require"
    if cfg.DisableTLS {
        sslMode = "disable"
    }

    q := make(url.Values)
    q.Set("sslMode", sslMode)
    q.Set("timezone", "utc")
    u := url.URL{
        Scheme: "postgres",
        User: url.UserPassword(cfg.User, cfg.Password),
        Host: cfg.Host,
        Path: cfg.Name,
        RawQuery: q.Encode(),
    }

    return sqlx.Open("postgres", u.String())
}

//StatusCheck returns nil if it can successfully talk to the database. It
//returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
    const q = "SELECT true"
    var tmp bool
    return db.QueryRowContext(ctx,q).Scan(&tmp)
}

func Log(query string, args ...interface{}) string {
    for i, arg := range args {
        n := fmt.Sprintf("$%d", i+1)
        var a string
        switch v := arg.(type) {
            case string:
                a = fmt.Sprintf("%q", v)
            case []byte:
                a = string(v)
            case []string:
                a = strings.Join(v, ",")
            default:
                a = fmt.Sprintf("%v", v)
        }

        query = strings.Replace(query, n, a, 1)
    }

    return query
}
