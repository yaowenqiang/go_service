package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "log"
    "time"
    "os"
    "fmt"
    "io/ioutil"

    "github.com/dgrijalva/jwt-go"
    "github.com/yaowenqiang/service/foundation/database"
    "github.com/yaowenqiang/service/business/data/schema"
)

/*
    openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
    openssl rsa -pubout -in private.pem -out public.pem
*/

func main() {
    //keygen()
    //tokengen()
    migrate()
}


func migrate() {
    dbConfig := database.Config{
        User: "postgres",
        Password: "postgres",
        Host: "0.0.0.0",
        Name: "postgres",
        DisableTLS: true,
    }

    db, err := database.Open(dbConfig)
    if err != nil {
        log.Fatalln(err)
    }
    defer db.Close()

    if err := schema.Migrate(db); err != nil {
        log.Fatalln(err)
    }

    if err := schema.Seed(db); err != nil {
        log.Fatalln(err)
    }

    fmt.Println("migrations complete")
}

func keygen() {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        log.Fatalln(err)
    }
    privateFile, err := os.Create("private.pem")
    if err != nil {
        log.Fatalln(err)
    }

    defer privateFile.Close()

    privateBlock := pem.Block{
        Type: "RSA PIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    }

    if err := pem.Encode(privateFile, &privateBlock); err != nil {
        log.Fatalln(err)
    }


    asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
    if err != nil {
        log.Fatalln(err)
    }

    //======================
    publicFile, err := os.Create("public.pem")
    if err != nil {
        log.Fatalln(err)
    }

    defer privateFile.Close()

    publicBlock := pem.Block{
        Type: "RSA PUBLIC KEY",
        Bytes: asn1Bytes,
    }


    if err := pem.Encode(publicFile, &publicBlock); err != nil {
        log.Fatalln(err)
    }


    fmt.Println("private and public key files generated")

}

func tokengen() {

    privatePEM, err := ioutil.ReadFile("private.pem")
    if err != nil {
        log.Fatalln(err)
    }


    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)

    if err != nil {
        log.Fatalln(err)
    }


    claims := struct {
        jwt.StandardClaims
        Roles []string `json:"roles"`
    }{
        StandardClaims: jwt.StandardClaims{
            Issuer: "service project",
            Subject: "123456789",
            ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
            IssuedAt: time.Now().Unix(),
        },
        Roles: []string{"ADMIN"},
    }

    method := jwt.GetSigningMethod("RS256")


    tkn := jwt.NewWithClaims(method, claims)

    tkn.Header["kid"] = "asdlfjldasjfdsjfldasjfl jlsjflweqjio;ewjejf"

    str, err := tkn.SignedString(privateKey)

    if err != nil {
        log.Fatalln(err)
    }


    fmt.Printf("----------------BEGIN TOKEN---------------\n%s\n----------------END TOKEN-----------------\n", str)


}

