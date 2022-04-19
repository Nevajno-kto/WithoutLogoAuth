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

const (
	SignUp = 0
	SignIn = 1
)

const AuthDelay int64 = 60

type AuthUseCase struct {
	authRepo    psql.AuthRepo
	clientsRepo psql.UsersRepo
	permRepo    psql.PemissionsRepo
}

var (
	DebugFileSignUp *os.File
	DebugFileSignIn *os.File
)

type DebugInfo struct {
	Phone string `json:"phone"`
	Code  int    `json:"code"`
}

func NewAuth(authRepo *psql.AuthRepo, usersRepo *psql.UsersRepo, permRepo *psql.PemissionsRepo) *AuthUseCase {
	//******************************** DEBUG ************************************
	DebugFileSignUp, _ = os.OpenFile("./signUpcode.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	DebugFileSignIn, _ = os.OpenFile("./signIncode.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	//******************************** DEBUG ************************************

	rand.Seed(time.Now().Unix())
	return &AuthUseCase{authRepo: *authRepo, clientsRepo: *usersRepo, permRepo: *permRepo}
}

func (uc *AuthUseCase) SignUp(ctx context.Context, u entity.Auth) (authjwt.Tokens, error) {
	var err error
	var user entity.User
	var tokens authjwt.Tokens

	if user, err = uc.clientsRepo.GetUser(ctx, u.User); err == nil {
		if user != (entity.User{}) {
			return tokens, fmt.Errorf("пользователь с таким номером телефона уже зарегистрирован")
		}
	} else {
		return tokens, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - SignUp: %w", err).Error())
	}

	switch u.Type {
	case "request":
		err = uc.RequestAuth(ctx, u, SignUp)
	case "accept":
		tokens, err = uc.AcceptSignUp(ctx, u)
	default:
		err = errors.Wrap(entity.ErrServiceProblem, "unknown value of type field")
	}

	return tokens, err
}

func (uc *AuthUseCase) SignIn(ctx context.Context, u entity.Auth) (authjwt.Tokens, error) {

	var err error
	var tokens authjwt.Tokens
	var user entity.User

	if user, err = uc.clientsRepo.GetUser(ctx, u.User); err == nil {
		if user == (entity.User{}) {
			return tokens, fmt.Errorf("пользователь с таким номером телефона не зарегистрирован")
		}
	} else {
		return tokens, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - SignIn: %w", err).Error())
	}

	switch u.Type {
	case "password":
		tokens, err = uc.SignInByPassword(ctx, u, user)
	case "request":
		err = uc.RequestAuth(ctx, u, SignIn)
	case "accept":
		tokens, err = uc.AcceptSignInByCode(ctx, u, user)
	default:
		err = errors.Wrap(entity.ErrServiceProblem, "unknown value of type field")
	}

	return tokens, err
}

func (uc *AuthUseCase) RequestAuth(ctx context.Context, u entity.Auth, sign int) error {
	var storeCode int
	var request_time int64
	var err error

	code := rand.Intn(8999) + 1000
	//_, err := sms.Send(u.User.Phone, fmt.Sprint(code))

	// if err != nil {
	// 	return fmt.Errorf("usecase - auth - SignUp - RequestSignUp - Send: %w", err)
	// }
	// defer f.Close()

	//******************************** DEBUG ************************************
	debugStr, _ := json.MarshalIndent(DebugInfo{Phone: u.User.Phone, Code: code}, " ", "")
	//f, _ := os.OpenFile("./code.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	//defer f.Close()
	if sign == 0 {
		if _, err = DebugFileSignUp.Write(debugStr); err != nil {
			panic(err)
		}
	} else {
		if _, err = DebugFileSignIn.Write(debugStr); err != nil {
			panic(err)
		}
	}
	//******************************** DEBUG ************************************

	if storeCode, request_time, err = uc.authRepo.GetAuthCode(ctx, u.User, sign); err == nil {
		if storeCode == 0 {
			err = uc.authRepo.InsertAuthCode(ctx, u.User, sign, code)
		} else {
			if request_time+AuthDelay > time.Now().Unix() {
				return fmt.Errorf(fmt.Sprintf("осталось %d секунд", time.Now().Unix()-(request_time+AuthDelay)))
			}
			err = uc.authRepo.UpdateAuthCode(ctx, u.User, sign, code)
		}
	}

	if err != nil {
		return errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - RequestAuth - GetAuthCode: %w", err).Error())
	}

	return nil
}

func (uc *AuthUseCase) AcceptSignUp(ctx context.Context, u entity.Auth) (authjwt.Tokens, error) {
	var storeCode int
	var request_time int64
	var err error
	var org string = u.User.Organization
	var permission entity.Permission

	storeCode, request_time, err = uc.authRepo.GetAuthCode(ctx, u.User, SignUp)

	if err != nil {
		return authjwt.Tokens{}, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - AcceptSignUp - GetAuthCode: %w", err).Error())
	}

	if storeCode != u.Code {
		return authjwt.Tokens{}, fmt.Errorf("неверный код подтверждения")
	}

	if request_time+3600 < time.Now().Unix() {
		return authjwt.Tokens{}, fmt.Errorf("время действия кода подтверждения вышло")
	}

	if u.User.Password != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(u.User.Password), 9)
		if err != nil {
			return authjwt.Tokens{}, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - AcceptSignUp - bcrypt.GenerateFromPassword %w", err).Error())
		}
		u.User.Password = string(pwd)
	}

	switch u.User.Type {
	case "client":
		permission, err = uc.permRepo.GetClientPermission(ctx, u.User)
	case "admin":
		permission, err = uc.permRepo.GetAdminPermission(ctx, u.User)
	default:
		return authjwt.Tokens{}, errors.Wrap(entity.ErrServiceProblem, "unknown value of user type field")
	}

	err = uc.clientsRepo.InsertUser(ctx, u.User)

	if err != nil {
		return authjwt.Tokens{}, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - AcceptSignUp - InsertUser %w", err).Error())
	}

	u.User, err = uc.clientsRepo.GetUser(ctx, u.User)
	u.User.Organization = org

	if err != nil {
		return authjwt.Tokens{}, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - AcceptSignUp - GetUser %w", err).Error())
	}

	err = uc.permRepo.InsertPermissionForUser(ctx, u.User, permission)

	if err != nil {
		return authjwt.Tokens{}, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - AcceptSignUp - InsertPermissionForUser %w", err).Error())
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authjwt.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(config.GetConfig().JWT.EatAuth))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: u.User.Id,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authjwt.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(config.GetConfig().JWT.EatRefresh))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: u.User.Id,
	})

	return authjwt.SidnedTokens(authToken, refreshToken, []byte(config.GetConfig().JWT.Secret))
}

func (uc *AuthUseCase) SignInByPassword(ctx context.Context, u entity.Auth, user entity.User) (authjwt.Tokens, error) {

	if u.User.Password == "" {
		return authjwt.Tokens{}, fmt.Errorf("пароль не заполнен")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.User.Password)); err != nil {
		return authjwt.Tokens{}, fmt.Errorf("введён не верный пароль")
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authjwt.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(config.GetConfig().JWT.EatAuth))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: user.Id,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authjwt.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(config.GetConfig().JWT.EatRefresh))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: user.Id,
	})

	return authjwt.SidnedTokens(authToken, refreshToken, []byte(config.GetConfig().JWT.Secret))
}

func (uc *AuthUseCase) AcceptSignInByCode(ctx context.Context, u entity.Auth, user entity.User) (authjwt.Tokens, error) {

	code, request_time, err := uc.authRepo.GetAuthCode(ctx, u.User, SignIn)

	if err != nil {
		return authjwt.Tokens{}, errors.Wrap(entity.ErrServiceProblem, fmt.Errorf("usecase - auth - AcceptSignInByCode - GetAuthCode: %w", err).Error())
	}

	if code != u.Code {
		return authjwt.Tokens{}, fmt.Errorf("неверный код подтверждения")
	}

	if request_time+3600 < time.Now().Unix() {
		return authjwt.Tokens{}, fmt.Errorf("время действия кода подтверждения вышло")
	}

	authToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authjwt.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(config.GetConfig().JWT.EatAuth))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: user.Id,
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authjwt.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(config.GetConfig().JWT.EatRefresh))),
			IssuedAt:  jwt.At(time.Now()),
		},
		UserId: user.Id,
	})

	return authjwt.SidnedTokens(authToken, refreshToken, []byte(config.GetConfig().JWT.Secret))
}
