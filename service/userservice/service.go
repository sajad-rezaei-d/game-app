package userservice

import (
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/phonenumber"
	"gameApp/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique (phoneNumber string) (bool,error)
	Register(u entity.User) (entity.User , error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User,bool,error)
	GetUserByID(userID uint) (entity.User,error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User)(string, error)
}

type Service struct{
	auth AuthGenerator
	repo Repository
}

func New(authGenerator AuthGenerator,repo Repository) Service{
	return Service{repo: repo, auth: authGenerator}
}

type RegisterRequest struct{
	Name string			`json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Password string		`json:"password"`
}

type UserInfo struct {
	ID uint					`json:"id"`
	PhoneNumber string		`json:"phone_number"`
	Name string				`json:"name"`
}

type RegisterResponse struct {
	User  UserInfo `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse,error) {
// TODO: we should verify phone number by verification code 
	//validate phone number
	if !phonenumber.IsValid(req.PhoneNumber){
		return  RegisterResponse{},fmt.Errorf("phone number is not valid")
	}
	//check uniqueness of phone number
	if isUnique ,err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err!=nil || !isUnique {
		if err!=nil {
			return RegisterResponse{},fmt.Errorf("unexpected error: %w",err)
		}

		if !isUnique {
			return RegisterResponse{},fmt.Errorf("phone number is not Unique")
		}
	} 
	
	
	//validate name 
	if len(req.Name)<3 {
		return RegisterResponse{},fmt.Errorf("name length should be greater than 3")
	}

	// TODO: check the password with regex pattern
	//validate password
	if len(req.Password)<8 {
		return RegisterResponse{},fmt.Errorf("password length should be greater than 8")
	}

	//create new user in storage 
	user := entity.User{
		ID: 0,
		PhoneNumber: req.PhoneNumber,
		Name: req.Name,
		Password: getBcryptHash(req.Password),
	}

	createdUser,err:=s.repo.Register(user)
	if err!=nil{
		return RegisterResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	//return created user
	return RegisterResponse{ UserInfo{
		ID: createdUser.ID,
		PhoneNumber: createdUser.PhoneNumber,
		Name: createdUser.Name,
	}},nil
}


func getBcryptHash(password string) string{
	pass := []byte(password)
	hPass,_ :=bcrypt.GenerateFromPassword(pass, 0)
	return string(hPass)
}


type LoginRequest struct{
	PhoneNumber string `json:"phone_number"`
	Password 	string `json:"password"`
}

type Tokens struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User UserInfo `json:"user"`
	Token Tokens  `json:"tokens"`
}

func (s Service) Login(req LoginRequest) (LoginResponse,error){
	const op= "userservice.login"
	// TODO: it would be better to user two separate method for existence check and getUserByPhoneNumber
	user,exist,err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err!=nil{
		return LoginResponse{},richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if !exist {
		return LoginResponse{},fmt.Errorf("record not found")
	}


	if bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.Password)) !=nil {
		return LoginResponse{},fmt.Errorf("username or password isn't correct")
	}

	// The first way (generate sessionID):
	// generate random session id
	// save sessionID in db for the given UserID
	// return session id to user

	// The second way (generate jwt token)
	accessToken, err := s.auth.CreateAccessToken(user)
	if err!=nil{
		return LoginResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	refreshToken ,err := s.auth.CreateRefreshToken(user)
	if err!=nil{
		return LoginResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	return LoginResponse{
		User : UserInfo{
			ID: user.ID,
			PhoneNumber: user.PhoneNumber,
			Name: user.Name,
		},
		Token: Tokens{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
		},
	},nil
}

type ProfileRequest struct {
	UserID uint	`json:"id"`
}

type ProfileResponse struct {
	Name string `json:"name"`
}

//all request inputs for interactor/service should be sanitized.

func (s Service) Profile(req ProfileRequest) (ProfileResponse,error) {
	const  op = "userservice.Profile"
	//getUserByID
	user, err := s.repo.GetUserByID(req.UserID)
	if err!=nil{
		return ProfileResponse{},richerror.New(op).WithErr(err).
		WithMeta(map[string]interface{}{"req": req})
		
	}

	return ProfileResponse{Name: user.Name},nil
}



