package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ishanjain28/pogo/common"
)

func RequireAuthorization() common.Handler {
	return func(rc *common.RouterContext, w http.ResponseWriter, r *http.Request) *common.HTTPError {
		if usr := decryptSession(r); usr != nil {
			rc.User = usr
			return nil
		}

		if strings.Contains(r.Header.Get("Accept"), "html") || r.Method == "GET" {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return nil
		} else {
			return &common.HTTPError{
				Message:    "Unauthorized!",
				StatusCode: http.StatusUnauthorized,
			}
		}

		return &common.HTTPError{
			Message:    "Unauthorized!",
			StatusCode: http.StatusUnauthorized,
		}
	}
}

func CreateSession(u *common.User, w http.ResponseWriter) error {

	// n_J6vaKjmmw4WB95DMorjQ.UMYdBLfttwPgQw9T0u0wdK7bGwDT9vwxoPAKWhjSAcpoiMsjh4eSfBkA4WB2deSoQu_cjCaJrcp77rvG67xkOeXsYpiclx2b-Oi7MHM3Kms.1507140277977.604800000.2CdxwiKAJT4SYJTVK-Du5jokr-CCnxo1ukdaVBkLRJg

	iv, err := generateRandomString(16)
	if err != nil {
		return err
	}
	userJSON, err := json.Marshal(u)
	if err != nil {
		return err

	}
	var hexedJSON []byte
	hex.Encode(hexedJSON, userJSON)

	fmt.Println(iv, string(userJSON), hexedJSON)

	block, err := aes.NewCipher(hexedJSON)
	if err != nil {
		return err
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	return nil
}

func decryptSession(r *http.Request) *common.User {

	c, err := r.Cookie("POGO_SESSION")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Printf("error in reading Cookie: %v", err)
		}
		return nil
	}
	fmt.Println(c)

	return nil
}

func generateRandomString(l int) ([]byte, error) {
	rBytes := make([]byte, l)

	_, err := rand.Read(rBytes)
	if err != nil {
		return nil, err
	}
	return rBytes, nil
}
