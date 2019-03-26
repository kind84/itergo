package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/dgrijalva/jwt-go"

	"github.com/julienschmidt/httprouter"
	"github.com/kind84/iterpro/repo"
	"golang.org/x/crypto/bcrypt"
)

var jk = []byte("my_super_secret_key")

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	r := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := repo.GetUser(r.Email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.Password))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ts, et, err := newJWT(u.Username, u.Role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    ts,
		Expires:  et,
		HttpOnly: true,
	})

	rr := struct {
		Token   string `json:"token"`
		Expires string `json:"expires"`
	}{
		Token:   ts,
		Expires: et.Format(time.RFC3339),
	}

	w = setHeaders(w)
	json.NewEncoder(w).Encode(&rr)
}

func Signup(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	r := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{}

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = nil
	u, err := repo.GetUser(r.Email)
	if err != nil && err != mgo.ErrNotFound {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err == nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "Username already taken")
		return
	}

	p, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u.Username = r.Email
	u.Password = string(p)
	u.Role = r.Role

	err = repo.SignupUser(&u)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ts, et, err := newJWT(u.Username, u.Role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    ts,
		Expires:  et,
		HttpOnly: true,
	})

	rr := struct {
		Token   string `json:"token"`
		Expires string `json:"expires"`
	}{
		Token:   ts,
		Expires: et.Format(time.RFC3339),
	}

	w = setHeaders(w)
	json.NewEncoder(w).Encode(&rr)
}

func newJWT(username, role string) (string, time.Time, error) {
	et := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: et.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString(jk)
	if err != nil {
		return "", time.Now(), err
	}
	return ts, et, nil
}

func authorize(role string, ts string) (bool, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(ts, claims, func(token *jwt.Token) (interface{}, error) {
		return jk, nil
	})

	if !t.Valid {
		return false, errors.New("Invalid token")
	}
	if err != nil {
		return false, err
	}

	u, err := repo.GetUser(claims.Username)
	if err != nil {
		return false, err
	}
	if u.Role != role {
		return false, errors.New("Invalid role")
	}

	return true, nil
}
