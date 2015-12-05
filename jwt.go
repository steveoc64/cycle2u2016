package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"reflect"
	"runtime"
	"time"
)

// location of the files used for signing and verification
const (
	privKeyPath = "keys/cycle2u.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/cycle2u.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

// keys are held in global variables
// i havn't seen a memory corruption/info leakage in go yet
// but maybe it's a better idea, just to store the public key in ram?
// and load the signKey on every signing request? depends on your usage i guess
var (
	verifyKey, signKey []byte
)

var SecurityRules map[string]interface{}

// read the JWT key files before starting http handlers
func _initJWT() {
	var err error

	signKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	verifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
	log.Println("... loaded RSA Keys")

	_initSecurityRules()
}

// Catch fatal conditions, and print a stack trace
func catchPanic(err *error, functionName string) {
	if r := recover(); r != nil {
		fmt.Printf("%s : PANIC Defered : %v\n", functionName, r)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	} else if err != nil && *err != nil {
		fmt.Printf("%s : ERROR : %v\n", functionName, *err)

		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		fmt.Printf("%s : Stack Trace : %s", functionName, string(buf))
	}
}

// Get the MD5 hash of a string, in HEX format
func GetMD5HexHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func generateToken(ID int, Role string, Username string) (string, error) {
	// create a signer for rsa 256
	t := jwt.New(jwt.GetSigningMethod("HS256"))

	// set our claims
	t.Claims["ID"] = ID
	t.Claims["Role"] = Role
	t.Claims["Username"] = Username

	// set the expire time for 30 days
	// see http://tools.ietf.org/html/draft-ietf-oauth-json-web-token-20#section-4.1.4
	t.Claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	tokenString, err := t.SignedString(signKey)
	return tokenString, err
}

//////////////////////////////////////////////////////////////////////////////////////////
//
func securityCheck(c *echo.Context, action string) (map[string]interface{}, error) {

	t := c.Request().Header["Authorization"]
	if len(t) < 1 {
		return nil, errors.New("No Auth Token")
	}
	if len(t) > 1 {
		return nil, errors.New("Too many Tokens")
	}

	token, err := jwt.Parse(t[0], func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid Token")
	}

	// Check the Role in the claim (inside the token) against the allowed roles for this page
	// Role can be passed in as either a single string, or a []string for multiple allowed roles
	claimedRole := token.Claims["Role"].(string)

	// Now get the matching roles for this action
	role := SecurityRules["defaultAllow"]
	if _, hasAction := SecurityRules[action]; hasAction {
		role = SecurityRules[action]
	} else {
		log.Println("Action", action, "not defined, so using the defaultAllow role instead = ", role)
	}

	v := reflect.TypeOf(role).Kind()
	switch v {
	case reflect.String:
		if role != "*" && claimedRole != role {
			return nil, errors.New("Invalid Role")
		}
	case reflect.Slice:
		ok := false
		for _, r := range role.([]string) {
			if r == claimedRole {
				ok = true
				break
			}
		}
		if !ok {
			return nil, errors.New("Invalid Role")
		}
	default:
		return nil, errors.New("Invalid Security Rule")
	}
	//	log.Println("Returning claims structure of :", token.Claims)
	return token.Claims, nil
}

func getClaimedUser(claim map[string]interface{}) (int, string) {
	UID := int(claim["ID"].(float64))
	Username := claim["Username"].(string)
	return UID, Username
}

func _initSecurityRules() {

	SecurityRules = make(map[string]interface{})

	SecurityRules["*"] = "*"
	SecurityRules["log"] = "*"
	SecurityRules["booking"] = []string{"public"}
	SecurityRules["admin"] = []string{"admin"}
}
