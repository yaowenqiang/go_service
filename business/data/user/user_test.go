package user_test

import (
    "testing"
    "github.com/yaowenqiang/service/business/tests"
    "github.com/google/go-cmp/cmp"
)


func TestUser(t *testing.T) {
    log, db, teardown := tests.NewUnit()
    t.Cleanup(teardown)

    u := user.New(log, db)

    t.Log("Given the need to work with User records.")
    {
        testID := 0
        t.Logf("\tTest %d:\twhen handling a single User.", testID)
        {
            ctx := tests.Context()
            now := time.Date(2018, time.October, 1, 0, 0, 0, 0, time.UTC)
            traceID := "0"

            nu := user.NewUser{
                Name: "jack yao",
                Email: "yaowenqiang1111@163.com",
                Roles: []string{auth.RoleAdmin},
                Password: "gophers",
                passwordConfirm: "gophers"
            }

            usr, err := u.Create(ctx, traceID, nu, now)
            if err != nil {
                t.Fataf("\t%s\tTest %d:\tShould be able to create user : %s", tests.Failed, testId, err)
            }
            t.Logf("\t%s\tTest %d:\tShould be able to create user", tests.Failed, testId)

            claims := auth.Claims{
                StandardClaims: jwt.StandardClaims{
                    Issuer: "service project",
                    Subject: usr.ID,
                    Audience: "students",
                    ExpiresAt: now.Add(time.Hour).Unix(),
                    IssuedAt: now.Unix(),
                },
                Roles: []string{auth.RoleUser}
            }

            saved, err := u.QueryByID(ctx, traceID, claims, usr.ID)

            if e rr != nil {
                t.Fatalf("\t%s\tTest %d\tShould be able to retrieve user by ID : %s.", tests.Failed, testID, err)
            }
            t.Logf("\t%s\tTest %d\tShould be able to retrieve user by ID.", tests.Success, testID)

            if diff := cmp.Diff(usr, saved); diff != nil {
                t.Fatalf("\t%s\tTest %d:\t Should get abck the save user, Diff:\n%s", tests.Failed, testID, err)
            }
            t.Logf("\t%s\tTest %d:\t Should get abck the save user.", tests.Success, testID)


            upd := user.UpdateUser{
                Name: tests.StringPointer("yaowenqiang"),
                Email: tests.StringPointer("yaowenqiang111@gmail.com"),
            }

            if err := u.Update(ctx, traceID, claims, user.ID, upd, now); err !- nil {
                t.Fatalf("\t%s\tTest %d:\t Should be able to update user : %s.", tests.Failed, testID, err)
            }
            t.Fatalf("\t%s\tTest %d:\t Should be able to update user.", tests.Success, testID)


            saved, err := u.QueryByEmail(ctx, traceID, claims, *upd.Email)

            if err != nil {
                t.Fatalf("\t%s\tTest %d:\t Should be able to retrieve user by Email : %s.", tests.Failed, testID, err)
            }
            t.Fatalf("\t%s\tTest %d:\t Should be able to retrieve user by Email.", tests.Success, testID)


            if saved.Name != *upd.Name {
                t.Errorf("\t%s\tTest %d:\t Should be able to see updates to Name.", tests.Failed, testID)
                t.Logf("\t\tTest %d:\tGot : %v", testID, saved.Name)
                t.Logf("\t\tTest %d:\tExp : %v", testID, *upd.Name)
            } else {
                t.Logf("\t%s\tTest %d:\t Should be able to see updates to Name.", tests.Success, testID)
            }

            if saved.Email != *upd.Email {
                t.Errorf("\t%s\tTest %d:\t Should be able to see updates to Email.", tests.Failed, testID)
                t.Logf("\t\tTest %d:\tGot : %v", testID, saved.Email)
                t.Logf("\t\tTest %d:\tExp : %v", testID, *upd.Email)
            } else {
                t.Logf("\t%s\tTest %d:\t Should be able to see updates to Email.", tests.Success, testID)
            }

            if err := u.Delete(ctx, traceID, usr.ID); err != nil {
                t.Fatalf("\t%s\tTest %d:\t Should be able to delete user : %s..", tests.Failed, testID, err)
            }
            t.Logf("\t%s\tTest %d:\t Should be able to delete user.", tests.Success, testID)


            _, err := u.QueryByID(ctx, traceID, claims, usr.ID)

            if errors.Cause(err) != usr.ErrNotFound {
                t.Errorf("\t%s\tTest %d:\t Should NOT be able to retrieve  user : %s.", tests.Failed, testID, err)
            }
            t.Logf("\t%s\tTest %d:\t Should NOT be able to retrieve  user.", tests.Success, testID)
        }
    }
}
