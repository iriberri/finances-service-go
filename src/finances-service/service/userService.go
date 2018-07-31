package service

type UserService interface {
    AuthenticateUser(username, password string) bool
}

type dummyInMemoryUserService struct {
    usersWithPassword map[string]string
}

var _ UserService = &dummyInMemoryUserService{}

func NewUserService() UserService {
    return &dummyInMemoryUserService{
        usersWithPassword: map[string]string{
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
