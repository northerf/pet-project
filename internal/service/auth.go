package service

import (
	"errors"
	"pet-project/internal/repository"
	"pet-project/pkg/model"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repository repository.UserRepository
	JwtSecret  []byte
}

func (s *AuthService) Register(email string, password string) error {
	_, err := s.Repository.FindByEmail(email)
	if err == nil {
		return errors.New("User already exists")
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		Email:    email,
		Password: string(hashedPwd),
		Name:     "NaN",
	}

	return s.Repository.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repository.FindByEmail(email)
	if err != nil {
		return "", errors.New("Email doesn't exists, user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.JwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) UpdateUser(user *model.User) error {
	return s.Repository.Update(user)
}

func (s *AuthService) DeleteUser(email string) error {
	user, err := s.Repository.FindByEmail(email)
	if err != nil {
		return err
	}
	return s.Repository.Delete(user)
}

func (s *AuthService) GetUserByID(id int) (*model.User, error) {
	return s.Repository.FindByID(id)
}

func (s *AuthService) DeleteUserByID(id int) error {
	user, err := s.Repository.FindByID(id)
	if err != nil {
		return err
	}
	return s.Repository.Delete(user)
}
