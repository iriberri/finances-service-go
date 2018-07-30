package service

import (
    "fmt"
    "sync"
    "testing"

    "github.com/stretchr/testify/assert"
)

type mockUserService struct{}

func (mockUserService) AuthenticateUser(username, password string) bool {
    return true
}

// Simulates charge on the `TokensService` to ensure thread safety.
func Test_TokensService_Charge(t *testing.T) {
    parallelExecs := 4
    triesPerExecs := 20000

    service := NewTokenService(&mockUserService{})
    wg := sync.WaitGroup{}
    wg.Add(parallelExecs * 2)

    startWg := sync.WaitGroup{}
    startWg.Add(1)

    for i := 0; i < parallelExecs; i++ {
        go func(prefix int) {
            startWg.Wait()
            for v := 0; v < triesPerExecs; v++ {
                token := fmt.Sprintf("token for max.mustermann.%d.%d", prefix, v)
                validated := service.ValidateToken(token)
                assert.False(t, validated, token)
            }
            wg.Done()
        }(i)
        go func(prefix int) {
            startWg.Wait()
            for v := 0; v < triesPerExecs; v++ {
                username := fmt.Sprintf("max.mustermann.%d.%d", prefix, v)
                password := ""

                token := service.CreateToken(username, password)
                assert.NotEmpty(t, token)

                recreatedToken := service.CreateToken(username, password)
                assert.Equal(t, token, recreatedToken)

                validated := service.ValidateToken(token)
                assert.True(t, validated)

                validated = service.ValidateToken(token)
                assert.True(t, validated)

                invalidated := service.InvalidateToken(token)
                assert.True(t, invalidated)

                invalidated = service.InvalidateToken(token)
                assert.False(t, invalidated)

                validated = service.ValidateToken(token)
                assert.False(t, validated)
            }
            wg.Done()
        }(i)
    }

    startWg.Done()
    wg.Wait()
}
