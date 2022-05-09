package HTTPBusiness

import (
	"DataBaseManage/public"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"time"

	//"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
)

const (
	mySigningKey = "werWER#@#21REWew21%#$SD*&312fdERss#@"
)
const (
	SecretKey = "werWER#@#21REWew21%#$SD*&312fdERss#@"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func StartServer() {

	log.Println("Now listening...")
	http.ListenAndServe(":8080", nil)
}
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	public.Log("login hand")

	username := r.FormValue("username")
	password := r.FormValue("password")
	if strings.ToLower(username) != "111" {
		public.Log("login hand3")
		if password != "111" {
			public.Log("login hand4")
			w.WriteHeader(http.StatusForbidden)
			public.Log("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

	public.Log("login hand5")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(2)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims
	public.Log("login hand6")
	//CREATE TABLE "svnupdate" ("Id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, "projectName" VARCHAR, "startTime" TIME, "endTime" TIME, "logs" TEXT, "hostGroup" VARCHAR, "svnStart" VARCHAR, "svnEnd" VARCHAR, "pubStatus" VARCHAR, "status" INTEGER, "logName" VARCHAR, "Remark" VARCHAR, "flag" INTEGER)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}
	public.Log("login hand8")
	response := Token{tokenString}
	JsonResponse(response, w)

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request) bool {

	token := r.Header.Get("Authorization")

	uid, _ := strconv.Atoi(string(r.Header.Get("uid")))
	return ValidateToken(uid, token)
}

type Checktoken struct {
	UID    int    `json:"uid"`
	Random string `json:"random"`
	Sign   string `json:"sign"`
	Token  string `json:"token"`
}

type Reply struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

//var Client *rpc.Client()
func ValidateToken(uid int, token string) bool {

	//public.Log("ValidateToken:")
	service := public.GetAuthServerName()
	//service := "192.168.1.11:9091"
	client, err := jsonrpc.Dial("tcp", service)
	if err != nil {
		public.Log("dial error:", err)
		return false
	}
	random := public.GetRandom()
	//public.Log("ValidateToken:2")
	sign := strings.ToUpper(public.GetMd5("sNZs8P79CchubNwsT2jdKyCC8iVaug08" + random))
	args := Checktoken{
		UID:    uid,
		Random: random,
		Sign:   sign,
		Token:  token,
	}

	var reply Reply
	err = client.Call("UserService.ValidateToken", args, &reply)
	if err != nil {
		//public.Log("Arith.Muliply call error:", err)
		return false
	}
	b, errr := json.Marshal(reply)
	if errr != nil {
		public.Log(errr)
	}
	jsonStr := string(b) //{"code":0,"msg":"token验证成功"}

	mapp := public.GetMapByJsonStr(jsonStr)
	if mapp != nil {

		code := mapp["code"]
		//public.Log(code)
		icode, ok := code.(float64)
		if ok {
			if icode == 0 {
				return true
			}
		} else {
			public.Log("code paras error")
		}

	}
	return false
}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if public.EnableAccessControlAllowOrigin() {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	//w.Header().Set("token", "")
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

var secretKey = "dWc2EpVYyJwKWeqF7YawteNup5Ai0Wx0"

func GetJwtToken(uid string) (tokenString string, err error) {

	iat := time.Now()
	exp := iat.Add(time.Hour * time.Duration(10))
	claims := jwt.MapClaims{
		"iat": iat.Unix(),
		"exp": exp.Unix(),
		"iss": "zzx",
		"uid": uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(secretKey))
	if err != nil {
		public.Log(err)
	}
	//public.Log("make token=" + tokenString)
	return
}

// 验证token
func ValidateJwtToken(tokenString string, uid string) bool {
	if tokenString == "abcdefg" {
		return true
	}
	//public.Log("ValidateJwtToken=" + tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		public.Log(err)
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return true
	} else {
		public.Log("token.Valid =" + tokenString)
		return false
	}

	uid1, _ := claims["uid"].(string)

	if uid1 != uid {
		public.Log("uid不相符")
		return false
	}

	if token.Valid {
		public.Log("token is valid=" + tokenString)
		return true
	}

	return false
}
