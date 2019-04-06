package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/julienschmidt/httprouter"
	"github.com/kind84/iterpro/repo"
	"golang.org/x/crypto/bcrypt"
)

var jk = []byte("my_super_secret_key")

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Refresh  bool   `json:"refresh"`
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

	rt, err := newRefresh(u.Username, u.Role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rr := struct {
		Token        string `json:"token"`
		Expires      string `json:"expires"`
		RefreshToken string `json:"refreshToken"`
	}{
		Token:        ts,
		Expires:      et.Format(time.RFC3339),
		RefreshToken: rt,
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
	if err != nil && err != mongo.ErrNoDocuments {
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

	rt, err := newRefresh(u.Username, u.Role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rr := struct {
		Token        string `json:"token"`
		Expires      string `json:"expires"`
		RefreshToken string `json:"refreshToken"`
	}{
		Token:        ts,
		Expires:      et.Format(time.RFC3339),
		RefreshToken: rt,
	}

	w = setHeaders(w)
	json.NewEncoder(w).Encode(&rr)
}

func RefreshToken(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	rt := getToken(req)
	if rt == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	c, err := authorize("", rt, true)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, err.Error())
		return
	}

	ts, et, err := newJWT(c.Username, c.Role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rt, err = newRefresh(c.Username, c.Role)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rr := struct {
		Token        string `json:"token"`
		Expires      string `json:"expires"`
		RefreshToken string `json:"refreshToken"`
	}{
		Token:        ts,
		Expires:      et.Format(time.RFC3339),
		RefreshToken: rt,
	}

	w = setHeaders(w)
	json.NewEncoder(w).Encode(&rr)
}

func newJWT(username, role string) (string, time.Time, error) {
	et := time.Now().Add(30 * time.Second)
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

func newRefresh(username, role string) (string, error) {
	et := time.Now().Add(365 * 24 * time.Hour)
	claims := &Claims{
		Username: username,
		Role:     role,
		Refresh:  true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: et.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString(jk)
	if err != nil {
		return "", err
	}
	return ts, nil
}

func authorize(role, ts string, refresh bool) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(ts, claims, func(t *jwt.Token) (interface{}, error) {
		return jk, nil
	})
	if err != nil {
		return claims, err
	}

	if claims, ok := token.Claims.(*Claims); !ok || claims.Refresh != refresh || !token.Valid {
		return claims, errors.New("Invalid token")
	}

	if role != "" && claims.Role != role {
		return claims, errors.New("Invalid role")
	}

	return claims, nil
}
