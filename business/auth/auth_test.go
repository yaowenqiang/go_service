package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/yaowenqiang/service/business/auth"
	"github.com/dgrijalva/jwt-go"
)


const (
    success = "\u2713"
    failed  = "\u2717"
)


func TestAuth(t *testing.T) {
    t.Log("Given the need to be able to authenticate and authorize access.")
    {
        testID := 0

        t.Logf("\t\tTest %d\t:when handling a single user.", testID)
        {
            privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
            if err != nil {
                log.Fatalln(err)
            }

            const KeyID = "sdjflkdasfjlkdasjfl;sdjfldsjfks"

            lookup := func(kid string) (*rsa.PublicKey, error) {
                switch kid {
                case KeyID:
                    return &privateKey.PublicKey, nil
                }

                return nil, fmt.Errorf("no publik key found for the specified kid: %s", kid)
            }


            a, err := auth.New("RS256", lookup, auth.Keys{KeyID: privateKey})

            if err != nil {
                t.Fatalf("\t%s\tTest %d\tShould be able to create an authenticator: %v", failed, testID, err)
            }
            t.Logf("\t%s\tTest %d:\tshould be able to parse the private key from pem", success, testID)


            claims := auth.Claims{
                StandardClaims: jwt.StandardClaims{
                    Issuer: "service project",
                    Subject: "sljfldsjfsdjfksdjfsdjf",
                    Audience: "students",
                    ExpiresAt: time.Now().Add(8750 * time.Hour).Unix(),
                    IssuedAt: time.Now().Unix(),
                },
                Roles: []string{auth.RoleAdmin},
            }

            token, err := a.GenerateToken(KeyID, claims)

            if err != nil {
                t.Fatalf("\t%s\ttest %d:\tShould be able to generate a JWT: %v", failed, testID, err)
            }
            t.Logf("\t%s\ttest %d:\tShould be able to generate a JWT:", success, testID)


            parsedClaims, err := a.ValidateToken(token)
            if err != nil {
                t.Fatalf("\t%s\ttest %d:\tShould be able to parse the claims: :", failed, testID)
            }

            t.Logf("\t%s\ttest %d:\tShould be able to parse the claims: :", success, testID)


            if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
                t.Logf("\t\ttest %d:\"exp: %d", testID, exp)
                t.Logf("\t\ttest %d:\"exp: %d", testID, got)
                t.Fatalf("\t%s\ttest %d:\tShould have the expected number of roles: %v:", failed, testID, err)
            }

            t.Logf("\t%s\ttest %d:\tShould have the expected number of roles:", success, testID)


            if exp, got := claims.Roles[0], parsedClaims.Roles[0]; exp != got {
                t.Logf("\t\ttest %d:\"exp: %s", testID, exp)
                t.Logf("\t\ttest %d:\"exp: %s", testID, got)
                t.Fatalf("\t%s\ttest %d:\tShould have the expected roles: %v:", failed, testID, err)

            }

            t.Logf("\t%s\ttest %d:\tShould have the expected roles:", failed, testID)
        }
    }
}
