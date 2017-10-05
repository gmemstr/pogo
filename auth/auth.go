package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ishanjain28/pogo/common"
)

const (
	enc = "cookie_session_encryption"
	// mac = "cookie_session_signature"

	// This is the key with which each cookie is encrypted, I'll recommend moving it to a env file.
	secret       = "super_long_string_difficult_to_crack"
	cookieName   = "POGO_SESSION"
	cookieExpiry = 60 * 60 * 24 * 30 // 30 days in seconds
)

func RequireAuthorization() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		usr, err := decryptCookie(r)
		if err != nil {
			fmt.Println(err.Error())
			if strings.Contains(r.Header.Get("Accept"), "html") || r.Method == "GET" {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return &common.HTTPError{
					Message:    "Unauthorized! Redirecting to /login",
					StatusCode: http.StatusTemporaryRedirect,
				}
			}
			return &common.HTTPError{
				Message:    "Unauthorized!",
				StatusCode: http.StatusUnauthorized,
			}
		}
		rc.User = usr
		return nil
	}
}

func CreateSession(u *common.User) (*http.Cookie, error) {

	iv, err := generateRandomString(16)
	if err != nil {
		return nil, err
	}
	userJSON, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	hexedJSON := hex.EncodeToString(userJSON)

	encKey := deriveKey(enc, secret)

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	// Fill the block with 0x0e
	if remBytes := len(hexedJSON) % aes.BlockSize; remBytes != 0 {
		t := []byte(hexedJSON)

		for i := 0; i < aes.BlockSize-remBytes; i++ {
			t = append(t, 0x0e)
		}
		hexedJSON = string(t)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	encCipher := make([]byte, len(hexedJSON)+aes.BlockSize)

	mode.CryptBlocks(encCipher, []byte(hexedJSON))

	cipherbase64 := base64urlencode(encCipher)
	ivbase64 := base64urlencode(iv)

	// Cookie format: iv.cipher.created_on.expire_on.HMAC
	cookieStr := fmt.Sprintf("%s.%s", ivbase64, cipherbase64)

	c := &http.Cookie{
		Name:   cookieName,
		Value:  cookieStr,
		MaxAge: cookieExpiry,
	}

	return c, nil
}

func decryptCookie(r *http.Request) (*common.User, error) {

	c, err := r.Cookie(cookieName)
	if err != nil {
		if err != http.ErrNoCookie {
			log.Printf("error in reading Cookie: %v", err)
		}
		return nil, err
	}

	csplit := strings.Split(c.Value, ".")
	if len(csplit) != 2 {
		return nil, errors.New("Invalid number of values in cookie")
	}

	ivb, cipherb := csplit[0], csplit[1]

	iv, err := base64urldecode(ivb)
	if err != nil {
		return nil, err
	}
	dcipher, err := base64urldecode(cipherb)
	if err != nil {
		return nil, err
	}

	if len(iv) != 16 {
		return nil, errors.New("IV length is not 16")
	}

	encKey := deriveKey(enc, secret)

	if len(dcipher)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext not multiple of blocksize")
	}

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, len(dcipher))

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(buf, []byte(dcipher))

	tstr := fmt.Sprintf("%x", buf)

	// Remove aes padding, 0e is used because it was used in encryption to mark padding
	padIndex := strings.Index(tstr, "0e")
	if padIndex == -1 {
		return nil, errors.New("Padding Index is -1")
	}
	tstr = tstr[:padIndex]

	data, err := hex.DecodeString(tstr)
	if err != nil {
		return nil, err
	}

	data, err = hex.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	u := &common.User{}
	err = json.Unmarshal(data, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func deriveKey(msg, secret string) []byte {
	key := []byte(secret)
	sha256hash := hmac.New(sha256.New, key)
	sha256hash.Write([]byte(msg))

	return sha256hash.Sum(nil)
}

func generateRandomString(l int) ([]byte, error) {
	rBytes := make([]byte, l)

	_, err := rand.Read(rBytes)
	if err != nil {
		return nil, err
	}
	return rBytes, nil
}

func base64urldecode(str string) ([]byte, error) {
	base64str := strings.Replace(string(str), "-", "+", -1)
	base64str = strings.Replace(base64str, "_", "/", -1)

	s, err := base64.RawStdEncoding.DecodeString(base64str)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func base64urlencode(str []byte) string {
	base64str := strings.Replace(string(str), "+", "-", -1)
	base64str = strings.Replace(base64str, "/", "_", -1)

	return base64.RawStdEncoding.EncodeToString([]byte(base64str))
}
