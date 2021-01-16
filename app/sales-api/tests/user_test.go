package tests

import (
    "net/http"
    "testing"
    "bytes"
    "strings"
    "encoding/json"
    "os"
    "net/http/httptest"

    "github.com/google/go-cmp/cmp"
	"github.com/yaowenqiang/service/business/tests"
	"github.com/yaowenqiang/service/app/sales-api/handlers"
	"github.com/yaowenqiang/service/business/data/user"
	"github.com/yaowenqiang/service/business/auth"
)

type UserTests struct {
    app http.Handler
    kid string
    userToken string
    adminToken string
}

func TestUsers(t *testing.T) {
    test := tests.NewIntergration(t)
    t.Cleanup(test.Teardown)

    shutdown := make(chan os.Signal, 1)

    tests := UserTests{
        app: handlers.API("develop", shutdown, test.Log, test.Auth, test.DB),
        kid: test.KID,
        userToken: test.Token(test.KID, "user@example.com", "gophers"),
        adminToken: test.Token(test.KID, "admin@example.com", "gophers"),
    }

    t.Run("crudUsers", tests.crudUser)

}

func (ut *UserTests) crudUser(t *testing.T) {
    nu := ut.postUser201(t)
    defer ut.deleteUser204(t, nu.ID)

    ut.getUser200(t, nu.ID)
    ut.putUser204(t, nu.ID)
    ut.putUser403(t, nu.ID)
}

func (ut *UserTests) postUser201(t *testing.T) user.Info {
    nu := user.NewUser{
        Name: "jack",
        Email: "jack@example.com",
        Roles: []string{auth.RoleAdmin},
        Password: "gophers",
        PasswordConfirm: "gophers",
    }

    body, err := json.Marshal(&nu)
    if err != nil {
        t.Fatal(err)
    }

    r := httptest.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    r.Header.Set("Authorization", "Bearer " + ut.adminToken)
    ut.app.ServeHTTP(w, r)

    var got user.Info

    t.Log("Given the need to create a new user with the users endpoint.")
    {
        testID := 0
        t.Logf("\tTest %d\tWhen using the declared user value.", testID)
        {
            if w.Code != http.StatusCreated {
                t.Fatalf("\t%s\tTest %d\tShould receive a status code of 201 for the response %v", tests.Failed, testID, w.Code)
            }

            t.Logf("\t%s\tTest %d\tShould receive a status code of 201 for the response",tests.Success, testID)

            if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
                t.Fatalf("\t%s\t %d\tShould be able to unMarshal the response :%v", tests.Failed, testID, err)
            }

            exp := got
            exp.Name = "jack"
            exp.Email = "jack@example.com"
            exp.Roles = []string{auth.RoleAdmin}

            if diff := cmp.Diff(got, exp); diff != "" {
                t.Fatalf("\t%s\tTest %d:\tShould get the expected result, Diff:\n%s", tests.Failed, testID, diff)
            }
            t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
        }

    }

    return got

}

func (ut *UserTests) deleteUser204(t *testing.T, id string) {

    r := httptest.NewRequest(http.MethodDelete, "/users/" + id, nil)
    w := httptest.NewRecorder()
    r.Header.Set("Authorization", "Bearer " + ut.adminToken)
    ut.app.ServeHTTP(w, r)

    t.Log("Given the need to validate deleting a user that does exist")
    {
        testID := 0
        t.Logf("\tTest %d\tWhen using the new  use %s. ", testID, id)
        {
            if w.Code != http.StatusNoContent {
                t.Fatalf("\t%s\tTest %d\tShould receive a status code of 204 for the response %v", tests.Failed, testID, w.Code)
            }

            t.Logf("\t%s\tTest %d\tShould receive a status code of 204 for the response",tests.Success, testID)

        }

    }

}

func (ut *UserTests) getUser200(t *testing.T, id string) {
    r := httptest.NewRequest(http.MethodGet, "/users/" + id, nil)
    w := httptest.NewRecorder()
    r.Header.Set("Authorization", "Bearer " + ut.adminToken)
    ut.app.ServeHTTP(w, r)

    t.Log("Given the need to validate getting a  user that exists")
    {
        testID := 0
        t.Logf("\tTest %d\tWhen using the new user %s", testID, id)
        {
            if w.Code != http.StatusOK {
                t.Fatalf("\t%s\tTest %d\tShould receive a status code of 200 for the response %v", tests.Failed, testID, w.Code)
            }

            t.Logf("\t%s\tTest %d\tShould receive a status code of 200 for the response",tests.Success, testID)


            var got user.Info

            if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
                t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
            }

            exp := got
            exp.ID = id
            exp.Name = "jack"
            exp.Email = "jack@example.com"
            exp.Roles = []string{auth.RoleAdmin}

            if diff := cmp.Diff(got, exp); diff != "" {
                t.Fatalf("\t%s\tTest %d:\tShould get the expected result, Diff:\n%s", tests.Failed, testID, diff)
            }
            t.Logf("\t%s\tTest %d:\tShould get the expected result.", tests.Success, testID)
        }

    }

}

func (ut *UserTests) putUser204(t *testing.T, id string) {

    body := `{"name":"jacky yao"}`

    r := httptest.NewRequest(http.MethodPut, "/users/" + id, strings.NewReader(body))
    w := httptest.NewRecorder()
    r.Header.Set("Authorization", "Bearer " + ut.adminToken)
    ut.app.ServeHTTP(w, r)

    t.Log("Given the need to update a user with the users endpoint.")
    {
        testID := 0
        t.Logf("\tTest %d:\tWhen using the modified user value.", testID)
        {
            if w.Code != http.StatusNoContent {
                t.Fatalf("\t%s\tTest %d\tShould receive a status code of 204 for the response %v", tests.Failed, testID, w.Code)
            }

            t.Logf("\t%s\tTest %d\tShould receive a status code of 204 for the response",tests.Success, testID)

            r = httptest.NewRequest(http.MethodGet, "/users/" + id, strings.NewReader(body))
            w = httptest.NewRecorder()
            r.Header.Set("Authorization", "Bearer " + ut.adminToken)
            ut.app.ServeHTTP(w, r)

            if w.Code != http.StatusOK {
                t.Fatalf("\t%s\tTest %d\tShould receive a status code of 200 for the retrieve %v", tests.Failed, testID, w.Code)
            }

            t.Logf("\t%s\tTest %d\tShould receive a status code of 200 for the retrieve", tests.Success, testID)

            var ru user.Info

            if err := json.NewDecoder(w.Body).Decode(&ru); err != nil {
                t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", tests.Failed, testID, err)
            }

            if ru.Name != "jacky yao" {
                t.Fatalf("\t%s\tTest %d:\tShould see an updated Name :  got %q want %q", tests.Failed, testID, ru.Name,"jacky yao" )
            }
            t.Logf("\t%s\tTest %d:\tShould see an updated Name.", tests.Success, testID)

            if ru.Email != "jack@example.com" {
                t.Fatalf("\t%s\tTest %d:\tShould not effect other fields like Email :  got %q want %q", tests.Failed, testID, ru.Email, "jack@example.com")
            }
            t.Logf("\t%s\tTest %d:\tShould not effect other fields like Email.", tests.Success, testID)

        }

    }

}

func (ut *UserTests) putUser403(t *testing.T, id string) {

    body := `{"name":"jacky yao"}`

    r := httptest.NewRequest(http.MethodPut, "/users/" + id, strings.NewReader(body))
    w := httptest.NewRecorder()
    r.Header.Set("Authorization", "Bearer " + ut.userToken)
    ut.app.ServeHTTP(w, r)

    t.Log("Given the need to update a user with the users endpoint.")
    {
        testID := 0
        t.Logf("\tTest %d:\tWhen using a non-admin user makes a request.", testID)
        {
            if w.Code != http.StatusForbidden {
                t.Fatalf("\t%s\tTest %d\tShould receive a status code of 403 for the response %v", tests.Failed, testID, w.Code)
            }
            t.Logf("\t%s\tTest %d\tShould receive a status code of 403 for the response",tests.Success, testID)
        }
    }

}
