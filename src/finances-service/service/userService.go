package service

import (
    "database/sql"

    log "github.com/sirupsen/logrus"
)

type UserService interface {
    AuthenticateUser(email, password string) bool
}

func NewUserService(databaseService DatabaseService) UserService {
    // todo #13: the configuration might ask for a dummy in-memory, for local dev and/or tests
    return &userService{
        databaseService: databaseService,
    }
}

//
// DATABASE-BASED IMPLEMENTATION
//

type userService struct {
    databaseService DatabaseService
}

var _ UserService = &userService{}

func (s *userService) AuthenticateUser(email, password string) bool {
    row := s.databaseService.QueryRow(
        "select id, display_name from users where email = $1",
        email)
    var id int
    var displayName string
    err := row.Scan(&id, &displayName)
    if err == sql.ErrNoRows {
        // no user with that email found
        return false
    }
    if err != nil {
        log.Infof("Failed to query for user. Refusing authentication. %s", err)
        return false
    }
    // todo #14 : Check the password
    // For the moment, this accepts the authentication the moment the email exists in the table.
    return true
}

//
// DUMMY IN-MEMORY IMPLEMENTATION
//

type dummyInMemoryUserService struct {
    usersWithPassword map[string]string
}

var _ UserService = &dummyInMemoryUserService{}

//noinspection GoUnusedFunction // todo #13: Offer per-configuration possibility to use this implementation
func newInMemoryDummyUserService() UserService {
    return &dummyInMemoryUserService{
        usersWithPassword: map[string]string{ // todo #13: Have this list configurable, not hard-coded
            "max.mustermann": "maxisthebest",
            "laura.gärtner":  "mysafepassword",
        },
    }
}

// Returns if the user authenticated with the provided username and password.
func (s dummyInMemoryUserService) AuthenticateUser(username, password string) (bool) {
    expectedPassword, found := s.usersWithPassword[username]
    if !found || password != expectedPassword {
        return false
    }
    return true
}
