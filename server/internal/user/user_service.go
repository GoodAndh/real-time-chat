package user

import (
	"context"
	"fmt"
	"realTime/config"
	"realTime/server/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repository Repository) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
	}
}

func (s *service) CreateUsers(c context.Context, req *RegisUserRequest) (*RegisUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	u := &User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
		CreatedAt:   time.Now(),
		LastUpdated: time.Now(),
	}

	uID, err := s.Repository.CreateUsers(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &RegisUserResponse{
		ID:       strconv.Itoa(uID),
		Username: u.Username,
		Email:    u.Email,
	}

	return res, nil

}

type JWTclaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	cfg      string
	jwt.RegisteredClaims
}

func (s *service) LoginUser(c context.Context, req *LoginUserRequest) (*LoginUserResponse, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("username atau password salah")
	}

	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		return nil, fmt.Errorf("username atau password salah")
	}

	jwtStx := JWTclaims{
		ID:               strconv.Itoa(int(user.ID)),
		Username:         user.Username,
		cfg:              config.InitConfig("main").JWTSecretKey,
		RegisteredClaims: jwt.RegisteredClaims{Issuer: strconv.Itoa(int(user.ID)), ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtStx)

	tokenString, err := token.SignedString([]byte(jwtStx.cfg))
	if err != nil {
		return nil, err
	}

	return &LoginUserResponse{
		accessToken: tokenString,
		ID:          (strconv.Itoa(int(user.ID))),
		Username:    user.Username,
	}, nil

}
