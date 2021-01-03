package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "log"
    "os"
    "fmt"
)

/*
    openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
    openssl rsa -pubout -in private.pem -out public.pem
*/

func main() {
    keygen()
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

