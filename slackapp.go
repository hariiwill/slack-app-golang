package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/mtmoses/httprouter"
)

/*
Server is the core structure for http router
*/
type Server struct {
	router *httprouter.Router
}

/*
Request is the core structure for web service input
*/
type Request struct {
	InputOne string `json:"inputone"`
	InputTwo string `json:"inputtwo"`
}

/*
Response is the core structure for web service output
*/
type Response struct {
	Status     bool    `json:"status"`
	Data       string  `json:"data,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
	Message    string  `json:"message"`
}

type SlackUser struct {
	Ok         bool      `json:"ok,omitempty"`
	UserMember []Members `json:"members,omitempty"`
}

type Members struct {
	Id      string  `json:"id,omitempty"`
	Profile Profile `json:"profile,omitempty"`
}

type Profile struct {
	Email string `json:"email,omitempty"`
}

type BlockInfo struct {
	Blocks  []blocks `json:"blocks,omitempty"`
	Channel string   `json:"channel,omitempty"`
}

type blocks struct {
	Type     string      `json:"type,omitempty"`
	Text     *text       `json:"text,omitempty"`
	Elements *[]elements `json:"elements,omitempty"`
}

type text struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type elements struct {
	Type string `json:"type,omitempty"`
	Text *text  `json:"text,omitempty"`
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Access-Control-Allow-Origin, Token, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Allow-Headers, *")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	s.router.ServeHTTP(w, r)
}

func showSplashscreen() {
	screenImage := `
	API HEALTHY
`
	fmt.Println(screenImage)
	fmt.Println("===============")
	fmt.Println("API")
	fmt.Println("===============")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintln(w, "HI the api is working")
}

func initializeRoutes() {
	port := "8060"
	url := "localhost"

	portString := ":" + port
	fmt.Println("Starting server on\n", url, portString)

	router := httprouter.New()
	router.GET("/", healthCheckHandler)
	router.POST("/user/v1/check", checkDegreeHandler)
	router.POST("/getdata", getDataHandler)

	http.ListenAndServe(":8060", &Server{router})
}

/* comprehentHandler()- Reads the input
 */
func checkDegreeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var names *Request

	//Reading request from the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Parsing Error", http.StatusInternalServerError)
	}

	//Unmarshaling the request to stay struct
	err = json.Unmarshal(body, &names)
	if err != nil {
		http.Error(w, "Parsing Error", http.StatusInternalServerError)
		return
	}

	res2D := Response{
		Status:  true,
		Message: "Successful",
		Data:    names.InputOne,
	}

	fmt.Fprintln(w, jSONResponse(res2D))
}

func main() {
	showSplashscreen()
	sendRemainderToUser()
	initializeRoutes()

}

func getDataHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var userInfo map[string]interface{}

	err := json.Unmarshal([]byte(r.FormValue("payload")), &userInfo)

	if err != nil {

		fmt.Println(err)
	}

	fmt.Println(userInfo)

}

func sendRemainderToUser() {

	token := "xoxb-960242739012-1130369187425-hSDF76DMYoehLoaxRpiU4fd4"

	baseURL := "https://slack.com/api/users.list"

	urlConstruct := fmt.Sprintf(baseURL + "?token=" + token)

	resp, err := http.Get(urlConstruct)

	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	var userinfo SlackUser

	err = json.Unmarshal(body, &userinfo)

	if err != nil {
		fmt.Println(err)
	}

	var blockdetail BlockInfo

	layout, _ := ioutil.ReadFile("slackrequest.json")
	json.Unmarshal([]byte(layout), &blockdetail)

	file, _ := json.MarshalIndent(blockdetail, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)

	for _, element := range userinfo.UserMember {

		blockdetail.Channel = element.Id

		s, _ := json.Marshal(blockdetail)

		resp := postRequest(s, "https://slack.com/api/chat.postMessage", token)

		fmt.Println(resp)

	}

}

func postRequest(s2 []uint8, url string, token string) []uint8 {

	client := &http.Client{}
	var bearer = "Bearer" + " " + token
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(s2))
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json; param=value")
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(f))
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	return f
}

func jSONResponse(resp Response) string {
	j, err := json.Marshal(&resp)
	if err != nil {
	}
	return string(j)
}
