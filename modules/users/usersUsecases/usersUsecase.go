package usersUsecases

import (
	"fmt"

	"github.com/TurterDev/basicgolang/config"
	"github.com/TurterDev/basicgolang/modules/users"
	"github.com/TurterDev/basicgolang/modules/users/usersRepositories"
	"github.com/TurterDev/basicgolang/pkg/basicgolangauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredentail) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetUserProfile(userId string) (*users.User, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	//Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	//insert User
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// InsertAdmin
func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {
	//Hashing password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	//insert User
	result, err := u.usersRepository.InsertUser(req, true)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// สำหรับ Sign-in | GetPassport คือการทำสำหรับเป็นบัตรผ่านว่าสามารถผ่านได้ไหม
func (u *usersUsecase) GetPassport(req *users.UserCredentail) (*users.UserPassport, error) {
	//Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	//Conmpare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	//Sign Token
	accessToken, err := basicgolangauth.NewBasicgolangAuth(basicgolangauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	//Refresh Token
	refreshToken, err := basicgolangauth.NewBasicgolangAuth(basicgolangauth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})

	//Set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	//Parse token
	claims, err := basicgolangauth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//Check Oauth
	oauth, err := u.usersRepository.FineOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	//Find Profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := basicgolangauth.NewBasicgolangAuth(
		basicgolangauth.Access,
		u.cfg.Jwt(),
		newClaims,
	)
	if err != nil {
		return nil, err
	}
	refreshToken := basicgolangauth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}
	if err := u.usersRepository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}


//Get User Profile
func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}