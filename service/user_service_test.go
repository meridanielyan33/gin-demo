package services_test

import (
	"context"
	"errors"
	"gin-demo/config"
	"gin-demo/middleware"
	"gin-demo/mocks"
	"gin-demo/model"
	"gin-demo/redis_utils"
	services "gin-demo/service"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister_Success(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	jwtStrategy := middleware.NewJWTStrategy(redis_utils.GetRedisClient())
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	user := &model.User{
		Username: "john",
		Email:    "john@example.com",
		Password: "secret",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

	err := svc.Register(user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegister_Failure(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	jwtStrategy := middleware.NewJWTStrategy(redis_utils.GetRedisClient())
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	user := &model.User{
		Username: "john",
		Email:    "john@example.com",
		Password: "secret",
	}

	mockRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(errors.New("failed to create user"))

	err := svc.Register(user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	config.InitTestConfig("testsecret")

	mockRepo := new(mocks.UserRepository)
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})
	defer rdb.FlushDB(context.Background())
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	mockRepo.On("FindByEmail", "john@example.com").
		Return(&model.User{Email: "john@example.com", Password: string(hashedPassword)}, nil)

	req := &model.UserLoginRequest{
		Email:    "john@example.com",
		Password: "secret",
	}

	resp, err := svc.Login(req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Contains(t, resp.Message, "Welcome")

	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	config.InitTestConfig("testsecret")

	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	mockRepo.On("FindByEmail", "john@example.com").
		Return(&model.User{Email: "john@example.com", Password: string(hashedPassword)}, nil)

	req := &model.UserLoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	resp, err := svc.Login(req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestLogin_RepoError(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	mockRepo.On("FindByEmail", "unknown@example.com").
		Return(nil, errors.New("not found"))

	req := &model.UserLoginRequest{Email: "unknown@example.com", Password: "whatever"}
	resp, err := svc.Login(req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "no such user")
}

func TestLogout_Success(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, mockRedis := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	user := &model.User{
		Username: "john",
		Email:    "john@example.com",
	}

	mockRepo.On("FindByEmail", "john@example.com").Return(user, nil)
	mockRedis.ExpectDel("john@example.com").SetVal(1)

	resp, err := svc.Logout(&model.UserLogoutRequest{Email: "john@example.com"})
	assert.NoError(t, err)
	assert.Contains(t, resp.Message, "logged out successfully")

	mockRepo.AssertExpectations(t)
	assert.NoError(t, mockRedis.ExpectationsWereMet())
}

func TestLogout_StrategyError(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, mockRedis := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	user := &model.User{Username: "john", Email: "john@example.com"}

	mockRepo.On("FindByEmail", "john@example.com").Return(user, nil)
	mockRedis.ExpectDel("john@example.com").SetErr(errors.New("redis unavailable"))

	resp, err := svc.Logout(&model.UserLogoutRequest{Email: "john@example.com"})

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "redis unavailable")

	mockRepo.AssertExpectations(t)
	assert.NoError(t, mockRedis.ExpectationsWereMet())
}

func TestLogout_UserNotFound(t *testing.T) {
	config.InitTestConfig("testsecret")

	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	mockRepo.On("FindByEmail", "missing@example.com").
		Return(nil, errors.New("not found"))

	req := &model.UserLogoutRequest{Email: "missing@example.com"}
	resp, err := svc.Logout(req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to find user")
}

func TestGetUsers(t *testing.T) {
	config.InitTestConfig("testsecret")

	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	expectedUsers := []model.UserData{
		{Username: "john", Email: "john@example.com"},
		{Username: "jane", Email: "jane@example.com"},
	}

	mockRepo.On("FindAll", "admin@example.com").Return(expectedUsers)

	users := svc.GetUsers("admin@example.com")

	require.Len(t, users, 2)
	assert.Equal(t, "john@example.com", users[0].Email)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers_Failure(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)
	mockRepo.On("FindAll", "admin@example.com").Return(nil)

	users := svc.GetUsers("admin@example.com")

	require.Nil(t, users)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)
	expectedUser := &model.UserData{Username: "john", Email: "john@example.com"}
	mockRepo.On("FindById", "123").Return(expectedUser, nil)

	user, err := svc.GetUserById("123")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "john@example.com", user.Email)
	mockRepo.AssertExpectations(t)
}

func TestGetUserById_LettersNotAllowed(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)
	mockRepo.On("FindById", "abc").Return(nil, errors.New("Invalid ID format"))
	user, err := svc.GetUserById("abc")
	require.Error(t, err)
	require.Nil(t, user)
	mockRepo.AssertExpectations(t)

}

func TestGetUserById_NotFound(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	mockRepo.On("FindById", "999").
		Return(nil, errors.New("user not found"))

	user, err := svc.GetUserById("999")

	require.Error(t, err)
	require.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")

	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	config.InitTestConfig("testsecret")
	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	expectedUser := &model.User{
		Username:  "john",
		Email:     "john@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	mockRepo.On("FindByEmail", "john@example.com").Return(expectedUser, nil)

	user, err := svc.GetUserByEmail("john@example.com")

	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.FirstName, user.FirstName)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	config.InitTestConfig("testsecret")

	mockRepo := new(mocks.UserRepository)
	rdb, _ := redismock.NewClientMock()
	jwtStrategy := middleware.NewJWTStrategy(rdb)
	svc := services.NewUserService(mockRepo, *jwtStrategy)

	mockRepo.On("FindByEmail", "missing@example.com").
		Return(nil, errors.New("user not found"))

	user, err := svc.GetUserByEmail("missing@example.com")

	require.Error(t, err)
	require.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")

	mockRepo.AssertExpectations(t)
}
