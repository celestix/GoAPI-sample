//   GoAPI-sample
//   Copyright (C) 2021 AnonyIndian (github.com/anonyindian)
//   @StarDevs 

//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU Affero General Public License for more details.


//Note: We have taken localhost as our endpoint in the following code's comments.

package main //our package

//Imported all the required packages
import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"unicode"
)

type server struct{}


//tokens which our API will validate.
var tokens []string = []string{
	"stardevs",
	"anonyindian",
	"golang-api",
	"3d23590d4eead1a56eb1b5cb490dec697a",
}





//Token generating function, required for returnToken function.
func genToken(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes) //returns our generated token as a string
}




//Our server's default Body
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status_code":"200","ok":true}`))
}




//function for localhost/check?.... 
//parameters: token, user_id
//example: localhost/check?token=stardevs&user_id=12383
func returnCheck(w http.ResponseWriter, r *http.Request) {
	var txt string
	var is_auth bool
	txt = ""
	desc := ""
	is_authorized, user_id := check_token(r) //validation of token and user_id extraction
	if is_authorized {
		is_auth = true
		if user_id == "" {
			desc = `,"error":"please provide a user_id"` //returns error if user_id is empty
		} else {
			for _, userID := range user_id {
				if unicode.IsDigit(userID) { 
					desc = fmt.Sprintf(`,"user_id":"%v"`, user_id) //returns on successfull request
				} else {
					desc = `,"error":"invalid user_id specified"` //returns error if user_id is not digits
				}
			}
		}
	} else {
		is_auth = false
		desc = `,"error":"please provide a valid api token"` //returns error if token is invalid
	}
	txt += fmt.Sprintf(`{"authorized":%v%v}`, is_auth, desc)
	w.Write([]byte(txt))
}





//function for localhost/token 
//example: localhost/token
//result: will return token in json formar {"ok":true,"token":"generated-token"}
func returnToken(w http.ResponseWriter, r *http.Request) {
	token := genToken(17) // Our generated token
	var txt string
	if token == "" {
		txt = `{"ok":false,"error":"an error occured while generating your api token"}` //returns error if token is empty
	} else {
		txt = fmt.Sprintf(`{"ok":true,"token":"%v"}`, token) //returns on successfull token generation
		tokens = append(tokens, token)
	}
	w.Write([]byte(txt))
}



//function to check if token is valid or not
//returns bool for validation and user_id as string
func check_token(r *http.Request) (bool, string) {
	is_authorized := false
	token := r.FormValue("token") //Form value of token arg

	for _, x := range tokens {
		if token == x {
			is_authorized = true //declares is_authorized = true when token is valid
		}
	}
	user_id := r.FormValue("user_id") //Form value of user_id arg

	return is_authorized, user_id
}



//Our main function where we'll point our http handler funcs
func main() {
	s := &server{}
	http.Handle("/", s)
	http.HandleFunc("/check", returnCheck)
	http.HandleFunc("/token", returnToken)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", 8080), nil)) //here 8080 is our port eg: localhost:8080
}
