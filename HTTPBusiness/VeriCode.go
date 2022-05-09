package HTTPBusiness

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/dchest/captcha"
)

var formTemplate = template.Must(template.New("example").Parse(formTemplateSrc))

func ShowVerifyCode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ShowVerifyCode")
	h := captcha.Server(captcha.StdWidth, captcha.StdHeight)
	h.ServeHTTP(w, r)
}
func GetVeriCode(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("GetVeriCode")
	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}

	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func Verify(captchaId, captchaValue string) bool {

	if captcha.VerifyString(captchaId, captchaValue) {
		return true
	} else {
		return false
	}
}
func VeriCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if !Verify(r.FormValue("captchaId"), r.FormValue("captchaValue")) {
		io.WriteString(w, "{\"code\":1,\"message\":\"wrong\",\"data\":\"\"}")
	} else {
		io.WriteString(w, "{\"code\":0,\"message\":\"success\",\"data\":\"\"}")
	}
	//io.WriteString(w, "<br><a href='/'>Try another one</a>")
}

var ShowVeriCode = captcha.Server(captcha.StdWidth, captcha.StdHeight)

/*
func main22() {
	http.HandleFunc("/", GetVeriCode)
	http.HandleFunc("/process", VeriCode)
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	fmt.Println("Server is at localhost:8666")
	if err := http.ListenAndServe("localhost:8666", nil); err != nil {
		log.Fatal(err)
	}
}
*/
const formTemplateSrc = `{
    "code": 0,
    "message": "success",
    "data": {
        "captchaId": "{{.CaptchaId}}",
        "imageUrl": "{{.CaptchaId}}.png"
    }
}
`
