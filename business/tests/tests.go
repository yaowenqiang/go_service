package tests
import (
    "testing"
    "log"
    "fmt"
    "os"
    "context"
    "time"
    "crypto/rand"
    "crypto/rsa"

	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
	"github.com/yaowenqiang/service/foundation/database"
	"github.com/yaowenqiang/service/foundation/web"
	"github.com/yaowenqiang/service/business/data/schema"
	"github.com/yaowenqiang/service/business/auth"
	"github.com/yaowenqiang/service/business/data/user"
)


const (
    Success = "\u2713"
    Failed = "\u2717"
)


var (
    dbImage = "postgres:13-alpine"
    dbPort = "5432"
    dbArgs = []string{"-e", "POSTGRES_PASSWORD=postgres"}
    UserID  = ""
    AdminID = ""
)

func NewUnit(t *testing.T) (*log.Logger, *sqlx.DB, func()) {
    c := startContainer(t, dbImage, dbPort, dbArgs...)
    cfg := database.Config{
        User: "postgres",
        Password: "postgres",
        Host: c.Host,
        Name: "postgres",
        DisableTLS: true,
    }
    db, err := database.Open(cfg)

    if err != nil {
        t.Fatalf("opening database connection: %v", err)
    }

    t.Log("waiting for database to be ready ...")

    var pingError error
    maxAttempts := 20
    for attempted := 1; attempted <= maxAttempts; attempted++ {
        pingError := db.Ping()
        if pingError == nil {
            break
        }
        time.Sleep(time.Duration(attempted) * 100 * time.Millisecond)
    }

    if pingError != nil {
        dumpContainerLogs(t, c.ID)
        stopContainer(t, c.ID)
        t.Fatalf("database never ready: %v", pingError)
    }

    if err := schema.Migrate(db); err != nil {
        stopContainer(t, c.ID)
        t.Fatalf("migrating error: %s", err)
    }

    teardown := func() {
        t.Helper()
        db.Close()
        stopContainer(t, c.ID)
    }

    log := log.New(os.Stdout, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
    return log, db, teardown
}

func Context() context.Context {
    values := web.Values{
        TraceID: uuid.New().String(),
        Now: time.Now(),
    }

    return context.WithValue(context.Background(), web.KeyValues, &values)
}

func StringPointer( s string ) *string {
    return &s
}

func IntPointer(i int) *int {
    return &i
}

type Test struct {
    TraceID string
    DB *sqlx.DB
    Log *log.Logger
    Auth *auth.Auth
    KID string

    t *testing.T
    Cleanup func()
}

func NewIntergration(t *testing.T) *Test {
    log, db, cleanup := NewUnit(t)

    if err := schema.Seed(db); err != nil {
        t.Fatal(err)
    }

    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

    if err != nil {
        t.Fatal(err)
    }

    kidID := ""

    lookup :=  func(kid string) (*rsa.PublicKey, error) {
        switch kid {
        case kidID:
            return &privateKey.PublicKey, nil
        }

        return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
    }

    auth, err := auth.New("RS256", lookup, auth.Keys{kidID: privateKey})

    if err != nil {
        t.Fatal(err)
    }

    test := Test{
        TraceID: "00000000-0000-0000-0000-000000000000",
        DB: db,
        Log: log,
        Auth: auth,
        KID: kidID,
        t: t,
        Cleanup: cleanup,
    }
    return &test;
}

func (test *Test) Teardown() {
    test.Cleanup()
}


func (test *Test) Token(kid string, email, pass string) string {
    u := user.New(test.Log, test.DB)
    claims, err := u.Authenticate(context.Background(), test.TraceID, time.Now(), email, pass)

    if err != nil {
        test.t.Fatal(err)
    }

    token, err := test.Auth.GenerateToken(kid, claims)

    if err != nil {
        test.t.Fatal(err)
    }

	//test.Log.Printf("generate token :%s", token)

    return token
}
