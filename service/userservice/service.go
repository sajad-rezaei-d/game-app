package userservice

import (
	"fmt"
	"gameApp/entity"
	"gameApp/pkg/phonenumber"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique (phoneNumber string) (bool,error)
	Register(u entity.User) (entity.User , error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User,bool,error)
	GetUserByID(userID uint) (entity.User,error)
}

type Service struct{
	signKey string
	repo Repository
}

func New(repo Repository,signKey string) Service{
	return Service{repo: repo, signKey: signKey}
}

type RegisterRequest struct{
	Name string			`json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Password string		`json:"password"`
}

type RegisterResponse struct {
	User entity.User
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
	return RegisterResponse{
		User:createdUser,
	},nil
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

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse,error){
	// TODO: it would be better to user two separate method for existence check and getUserByPhoneNumber
	user,exist,err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err!=nil{
		return LoginResponse{},fmt.Errorf("unexpected error: %w",err)
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
	token, err :=createToken(user.ID, s.signKey)
	if err!=nil{
		return LoginResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	return LoginResponse{AccessToken: token},nil
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

//all request inputs for interactor/service should be sanitized.

func (s Service) Profile(req ProfileRequest) (ProfileResponse,error) {
	//getUserByID
	user, err := s.repo.GetUserByID(req.UserID)
	if err!=nil{
		// I have not expect the repository call return "record not found" error,
		//  because I assume the interactor input is sanitized.
		// TODO: we can use Rich Error.
		return ProfileResponse{},fmt.Errorf("unexpected error: %w",err)
	}

	return ProfileResponse{Name: user.Name},nil
}



type Claims struct {
	jwt.RegisteredClaims
	UserID uint
}
//TODO: use method RS256 
func createToken(userID uint, signKey string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		UserID: userID,
	})

	return t.SignedString([]byte(signKey))
}