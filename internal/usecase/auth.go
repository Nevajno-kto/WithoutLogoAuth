package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nevajno-kto/without-logo-auth/config"
	"github.com/nevajno-kto/without-logo-auth/internal/entity"
	"github.com/nevajno-kto/without-logo-auth/internal/usecase/repo/psql"
	authjwt "github.com/nevajno-kto/without-logo-auth/pkg/jwt"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type ClientRepo interface {
	Get()
	Insert()
}

type AuthUseCase struct {
	authRepo    psql.AuthRepo
	clientsRepo psql.UsersRepo
}

var (
	DebugFileSignUp *os.File
	DebugFileSignIn *os.File
)

type DebugInfo struct {
	Phone string `json:"phone"`
	Code  int    `json:"code"`
}

func NewAuth(authRepo *psql.AuthRepo, clientsRepo *psql.UsersRepo) *AuthUseCase {
	//******************************** DEBUG ************************************
	DebugFileSignUp, _ = os.OpenFile("./signUpcode.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	DebugFileSignIn, _ = os.OpenFile("./signUpcode.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	//******************************** DEBUG ************************************

	rand.Seed(time.Now().Unix())
	return &AuthUseCase{authRepo: *authRepo, clientsRepo: *clientsRepo}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, u entity.Auth) error {
	var err error

	if client, err := uc.clientsRepo.GetUser(ctx, u.User); err == nil {
		if client != (entity.User{}) {
			return errors.Wrap(entity.ErrSignUp, "пользователь с таким номером телефона уже зарегистрирован")
		}
	} else {
		return fmt.Errorf("usecase - auth - SignUp: %w", err)
	}

	switch u.Type {
	case "request":
		err = uc.RequestSignUp(ctx, u)
	case "accept":
		err = uc.AcceptSignUp(ctx, u)
	default:
		err = errors.New("unknown value of type field")
	}

	return err
}

func (uc *AuthUseCase) RequestSignUp(ctx context.Context, u entity.Auth) error {

	code := rand.Intn(899999) + 100000
	//_, err := sms.Send(u.User.Phone, fmt.Sprint(code))
	//

	// if err != nil {
	// 	return fmt.Errorf("usecase - auth - SignUp - RequestSignUp - Send: %w", err)
	// }
	// defer f.Close()

	//******************************** DEBUG ************************************
	debugStr, _ := json.MarshalIndent(DebugInfo{Phone: u.User.Phone, Code: code}, " ", "")
	//f, _ := os.OpenFile("./code.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	//defer f.Close()

	if _, err := DebugFileSignUp.Write(debugStr); err != nil {
		panic(err)
	}
	//******************************** DEBUG ************************************

	var storeCode int
	var err error

	if storeCode, _, err = uc.authRepo.GetSignUpCode(ctx, u.User); err == nil {
		if storeCode == 0 {
			err = uc.authRepo.InsertSignUpCode(ctx, u.User, code)
		} else {
			err = uc.authRepo.UpdateSignUpCode(ctx, u.User, code)
		}
	}

	if err != nil {
		return fmt.Errorf("usecase - auth - SignUp - RequestSignUp - GetSignUpCode: %w", err)
	}

	return nil
}

func (uc *AuthUseCase) AcceptSignUp(ctx context.Context, u entity.Auth) error {

	code, request_time, err := uc.authRepo.GetSignUpCode(ctx, u.User)

	if err != nil {
		return fmt.Errorf("usecase - auth - SignUp - AcceptSignUp: %w", err)
	}

	if code != u.Code {
		return errors.Wrap(entity.ErrSignUp, "неверный код подтверждения")
	}

	if request_time+3600 < time.Now().Unix() {
		return errors.Wrap(entity.ErrSignUp, "время действия кода подтверждения вышло")
	}

	if u.User.Password != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(u.User.Password), 9)
		if err != nil {
			return fmt.Errorf("usecase - auth - SignUp - AcceptSignUp - bcrypt.GenerateFromPassword %w", err)
		}
		u.User.Password = string(pwd)
	}

	err = uc.clientsRepo.InsertUser(ctx, u.User)

	if err != nil {
		return errors.Wrap(entity.ErrSignUp, "не удалось зарегистрировать пользователя")
	}

	return nil
}

func (uc *AuthUseCase) SignIn(ctx context.Context, u entity.Auth) (string, error) {

	var err error
	var token string
	var user entity.User

	if user, err = uc.clientsRepo.GetUser(ctx, u.User); err == nil {
		fmt.Println(user.Id)
		if user == (entity.User{}) {
			return token, errors.Wrap(entity.ErrSingIn, "пользователь с таким номером телефона не зарегистрирован")
		}
	} else {
		return token, fmt.Errorf("usecase - auth - SignIn: %w", err)
	}

	switch u.Type {
	case "password":
		token, err = uc.SignInByPassword(ctx, u, user)
	case "request":
		token, err = uc.RequestSignInByCode(ctx, u)
	case "accept":
		token, err = uc.AcceptSignInByCode(ctx, u)
	default:
		err = errors.New("unknown value of type field")
	}

	return token, err
}

func (uc *AuthUseCase) SignInByPassword(ctx context.Context, u entity.Auth, user entity.User) (string, error) {

	if u.User.Password == "" {
		return "", errors.Wrap(entity.ErrSingIn, "пароль не заполнен")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.User.Password)); err != nil {
		return "", errors.Wrap(entity.ErrSingIn, "введён не верный пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, authjwt.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(config.GetConfig().JWT.Eat))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: user.Id,
	})

	return token.SignedString([]byte(config.GetConfig().JWT.Secret))
}

func (uc *AuthUseCase) RequestSignInByCode(ctx context.Context, u entity.Auth) (string, error) {

	code := rand.Intn(899999) + 100000

	//******************************** DEBUG ************************************
	debugStr, _ := json.MarshalIndent(DebugInfo{Phone: u.User.Phone, Code: code}, " ", "")
	//f, _ := os.OpenFile("./code.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	//defer f.Close()

	if _, err := DebugFileSignUp.Write(debugStr); err != nil {
		panic(err)
	}
	//******************************** DEBUG ************************************

	return "", nil
}

func (uc *AuthUseCase) AcceptSignInByCode(ctx context.Context, u entity.Auth) (string, error) {

	return "", nil
}
