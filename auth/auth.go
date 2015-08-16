package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jkusniar/lara/app"
	"golang.org/x/crypto/scrypt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	SALT_BYTES    = 32
	HASH_BYTES    = 64
	PRIV_KEY_PATH = "lara.rsa"     // openssl genrsa -out lara.rsa 2048
	PUB_KEY_PATH  = "lara.rsa.pub" // openssl rsa -in lara.rsa -pubout > lara.rsa.pub
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey // TODO better to load PrivateKey only when signing (login process)
)

func init() {
	// TODO panic should be fatal (os.exit(1))
	signBytes, err := ioutil.ReadFile(PRIV_KEY_PATH)
	if err != nil {
		panic(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	verifyBytes, err := ioutil.ReadFile(PUB_KEY_PATH)
	if err != nil {
		panic(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		panic(err)
	}
}

type loginMsg struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

// Login performs login process
// it parses json content of POST http request, validates users password, and generates JWT token
// token is returned in response body
func Login(w http.ResponseWriter, r *http.Request) {
	var m loginMsg
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.Log.Errorf("Error parsing login request body: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO authentication
	// 1. load salt from DB based on username. Error if user not exists
	// 2. hash password from json message
	// 3. compare hashes
	// 4. load user permissions
	// 5. generate JWT with username and permissions

	// generate JWT
	t := jwt.New(jwt.SigningMethodRS256)
	// set our claims
	t.Claims["UserInfo"] = struct {
		Name string
		Role string
	}{m.User, "user"}
	// set the expire time
	t.Claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := t.SignedString(signKey)
	if err != nil {
		app.Log.Errorf("Token Signing error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, tokenString)
}

// Register registeres new user
// creates salt, password hash and stores them in database with username
func Register(w http.ResponseWriter, r *http.Request) {
	// decode request
	var m loginMsg
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.Log.Errorf("Error parsing register user request body: %s\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// generate salt
	salt := make([]byte, SALT_BYTES)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		app.Log.Errorf("Error generating salt: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	app.Log.Debugf("Generated salt: %x\n", salt)

	// generate hash
	var hash []byte
	hash, err = hashPass(m.Pass, salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO store username, hash and salt in DB
	w.WriteHeader(200)
	app.Log.Debug(hash)
}

// Validate validates if request containst proper authorization header with valid authorization token
// returns error if validation fails, otherwise returns nil
func Validate(w http.ResponseWriter, r *http.Request) error {
	ah := r.Header.Get("Authorization")
	if ah == "" {
		return errors.New("Authorization header missing")
	}

	if len(ah) <= 7 || strings.ToUpper(ah[0:6]) != "BEARER" {
		return errors.New("Authorization header shoud contain \"Bearer\" <JWT token>")
	}

	tokenString := ah[7:]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// required because of JWT spec vulnerability
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// might lookup different key based on user (client) info in token
		// so far implementation supports only one keypair
		return verifyKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid authorization token")
	}

	app.Log.Debug("Granted access to restricted area")
	return nil
}

func hashPass(password string, salt []byte) (hash []byte, err error) {
	hash, err = scrypt.Key([]byte(password), salt, 16384, 8, 1, HASH_BYTES)
	if err != nil {
		app.Log.Errorf("Scrypt failed: %s\n", err)
		return
	}

	app.Log.Debugf("Hash of password %s is %x\n", password, hash)
	return
}
