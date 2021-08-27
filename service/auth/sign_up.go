package auth

import (
	"context"
	"errors"
	desc "github.com/auth_service/api"
	"github.com/auth_service/models"
	"github.com/auth_service/tools"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
	"unicode"
)

const (
	bcryptCost = 11
)

func (s *Service) SignUp(ctx context.Context, req *desc.SignUpRequest) (err error) {
	const (
		checkEmailExistence = `select exists(select * from user_data where email = $1);`
		dbRequestToAddData  = `insert  into user_data (name, last_name, email, pass_hash) values 
			($1, $2, $3, $4) returning id`
	)
	var (
		passHash    []byte
		emailExists bool
		userID      int32
	)

	if err = checkThePass(req.Pass); err != nil {
		return
	}
	if err = checkTheUserData(req.Name, req.LastName); err != nil {
		return
	}

	tx, err := s.db.Begin()
	if err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при создании транзакции", http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			er := tx.Rollback()
			if er != nil {
				log.Printf("error while rolling up the transaction: %v", er)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			err = tools.NewErrorMessage(err, "Ошибка при tx.Commit", http.StatusInternalServerError)
		}
	}()

	if req.Name == "" || req.LastName == "" || req.Pass == "" {
		err = tools.NewErrorMessage(errors.New("bad request"), "Какое то из полей пустое", http.StatusBadRequest)
		return
	}

	if err = tx.QueryRowEx(ctx, checkEmailExistence, nil, req.PhoneNumber).Scan(&emailExists); err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при проверке на сущестование емейла", http.StatusInternalServerError)
		return
	}

	if emailExists {
		err = tools.NewErrorMessage(errors.New("this email is already registered"),
			"Данный емейл уже зарегестрирован", http.StatusBadRequest)
		return
	}

	if passHash, err = bcrypt.GenerateFromPassword([]byte(req.Pass), bcryptCost); err != nil {
		err = tools.NewErrorMessage(err, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}

	if err = tx.QueryRowEx(ctx, dbRequestToAddData, nil, req.Name, req.LastName, req.PhoneNumber, string(passHash)).Scan(&userID); err != nil {
		err = tools.NewErrorMessage(err, "Ошибка при сохранении данных", http.StatusInternalServerError)
		return
	}
	//
	//output.AccessToken, err = generateToken(userID)
	//if err != nil {
	//	err = tools.NewErrorMessage(err, "Ошибка при создании токена", http.StatusInternalServerError)
	//}
	return
}

func checkThePass(pass string) (err error) {
	if len(pass) != len([]rune(pass)) {
		err = tools.NewErrorMessage(errors.New("bad pass"),
			"Некорректный пароль",
			http.StatusBadRequest)
		return
	}

	if len(pass) < 8 || len(pass) > 50 {
		err = tools.NewErrorMessage(errors.New("bad pass len"),
			"Длина пароля должна быть от 8 до 50 символов", http.StatusBadRequest)
	}
	return
}

func checkTheUserData(firstName, lastName string) (err error) {
	fieldsArr := []string{firstName, lastName}

	for i := range fieldsArr {
		fieldRune := []rune(fieldsArr[i])
		for i := range fieldRune {
			if !unicode.IsLetter(fieldRune[i]) {
				err = tools.NewErrorMessage(errors.New("bad field"),
					"Следует использовать только символы", http.StatusBadRequest)
				return
			}
		}
	}
	return
}

func generateToken(userID int32) (response string, err error) {
	claims := models.ClaimWithID{
		ID: strconv.Itoa(int(userID)),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 3600,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	response, err = token.SignedString(models.JwtSigningKey)
	return
}
