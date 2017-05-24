package main

import (
	_ "bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/securecookie"
	"golang.org/x/net/context"
)

const (
	defaultformat    = "01/02/2006"
	retriveCondition = "and description not like 'Transfer%' and description not like 'MOB PAYMENT%'"
)

var startInterval time.Time
var Last_update time.Time

type AuthorizeInfo struct {
	username        string
	authorizationId map[string]int
	authorizedUntil time.Time
}
type Service struct{}
type RetrieveRequest struct {
	RE []map[string]interface{} `json:"response"`
}
type Response struct {
	RE interface{} `json:"response"`
}
type DefaultResponse struct {
	RE map[string]interface{} `json:"response"`
}
type RetrieveParameters struct {
	dateto   string
	datefrom string
}
type login struct {
	username string
	password string
}
type USERS struct {
	salt     string
	password string
}
type GenError struct {
	s         string
	Httperror int
}
type ResponseWithCookie struct {
	Response string
	username string
}
type pass struct {
	cookie   *http.Cookie
	Retrieve interface{}
}

func (e GenError) Error() string {
	return e.s
}
func (e GenError) getHttpError() int {
	return e.Httperror
}
func (Service) CheckParameters(request RetrieveParameters) bool {
	re := regexp.MustCompile(`^[0-9]{2,2}\/[0-9]{2,2}\/([0-9]{4,4}|[0-9]{2,2})$`)
	matched := re.MatchString(request.dateto)
	if !matched {
		return false
	}
	matched = re.MatchString(request.datefrom)
	if !matched {
		return false
	}
	return true
}
func (s Service) Retrieve(request RetrieveParameters) ([]map[string]interface{}, error) {
	ok := s.CheckParameters(request)
	if !ok {
		return nil, errors.New("Parameters assigned to blank are invalid")
	}
	condition := fmt.Sprintf("WHERE DATE BETWEEN DATE '%s' and DATE '%s'", request.datefrom, request.dateto)
	req, _ := SQL_exec("SELECT * FROM TRANSACTIONS "+condition+retriveCondition, nil)
	//log.Println(req)
	return req, nil
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(cookie *http.Cookie) (username string) {
	cookieValue := make(map[string]string)
	if err := cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
		username = cookieValue["name"]
	}
	return username
}
func setSession(userName string) *http.Cookie {
	t := RandStringBytesMaskImprSrc(64)
	value := map[string]string{
		"name":    userName,
		"session": t,
	}
	if user, ok := usersAllowed[userName]; !ok {
		usersAllowed[userName] = AuthorizeInfo{username: userName, authorizationId: map[string]int{t: 1}, authorizedUntil: time.Now().Add(15 * time.Minute)}
	} else {
		user.authorizationId[t] = 1
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:     "session",
			Value:    encoded,
			Path:     "/",
			MaxAge:   60 * 15,
			HttpOnly: true,
		}
		return cookie
	}
	return clearSession()
}

func clearSession() *http.Cookie {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	return cookie
	//http.SetCookie(response, cookie)
}

var usersAllowed map[string]AuthorizeInfo

func main() {
	loc, _ := time.LoadLocation("UTC")
	startInterval, _ = time.ParseInLocation("2006-01-02", "2016-01-04", loc)
	usersAllowed = make(map[string]AuthorizeInfo, 0)
	ctx := context.Background()
	svc := Service{}
	Last_update = time.Time{}
	Initilize()
	StartDispatcher(1)
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}
	//var RE endpoint.Endpoint
	//RE = RetrieveEndpoint(svc)
	//RE = Authorize()(RE)
	Retrievehandler := kithttp.NewServer(
		ctx,
		Authorize()(RetrieveEndpoint(svc)),
		decodeRetrieveRequest,
		encodeResponse,
		opts...,
	)
	Updatehandler := kithttp.NewServer(
		ctx,
		Authorize()(UpdateEndpoint(svc)),
		func(_ context.Context, r *http.Request) (interface{}, error) {
			cookieUser, err := r.Cookie("session")
			if err != nil {
				return nil, GenError{"No Current Session is good aaa", http.StatusUnauthorized}
			}
			return pass{cookieUser, nil}, nil
			//return pass{cookieUser, RetrieveParameters{request.Get("dateto"), request.Get("datefrom")}}, nil
		},
		encodeResponse,
		opts...,
	)
	Loginhandler := kithttp.NewServer(
		ctx,
		LoginEndpoint(svc),
		decodeLoginRequest,
		encodeResponseWithCookie,
		opts...,
	)
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 7 * time.Second,
		Addr:         ":40500",
	}
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `<!DOCTYPE html>
					<!--[if lt IE 7]> <html class="lt-ie9 lt-ie8 lt-ie7" lang="en"> <![endif]-->
					<!--[if IE 7]> <html class="lt-ie9 lt-ie8" lang="en"> <![endif]-->
					<!--[if IE 8]> <html class="lt-ie9" lang="en"> <![endif]-->
					<!--[if gt IE 8]><!--> <html lang="en"> <!--<![endif]-->
					<head>
						<meta charset="utf-8">
						<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
						<title>Login Form</title>
						<link rel="stylesheet" href="css/style.css">
						<!--[if lt IE 9]><script src="//html5shim.googlecode.com/svn/trunk/html5.js"></script><![endif]-->
					</head>
					<body>
						<section class="container">
							<div class="login">
								<h1>Login to Web App</h1>
								<form method="post" action="login">
									<p><input type="text" name="user" value="" placeholder="Username or Email"></p>
									<p><input type="password" name="password" value="" placeholder="Password"></p>
									<p class="remember_me">
										<label>
											<input type="checkbox" name="remember_me" id="remember_me">
											Remember me on this computer
										</label>
									</p>
									<p class="submit"><input type="submit" name="commit" value="Login"></p>
								</form>
							</div>

						</section>
					</body>
					</html>`)
		return
	})
	http.Handle("/Retrieve", Retrievehandler)
	http.Handle("/Update", Updatehandler)
	http.Handle("/login", Loginhandler)
	//http.HandleFunc("/Create", create)
	http.HandleFunc("/Interval", timeInterval)
	http.HandleFunc("/Logout", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "https://willshadow.dynu.net:40500/home", 302)
			return
		}
		http.SetCookie(w, clearSession())
		http.Redirect(w, r, "https://willshadow.dynu.net:40500/home", 302)
	})
	log.Println("Starting server on 40500")
	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}
func timeInterval(w http.ResponseWriter, r *http.Request) {
	rightNow := time.Now()
	var startTime = startInterval
	for startTime.Before(rightNow) {
		startTime = startTime.AddDate(0, 0, 28)
	}
	var endTime = startTime
	startTime = startTime.AddDate(0, 0, -28)
	endTime = rightNow
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	e.Encode(DefaultResponse{map[string]interface{}{"datefrom": startTime.Format(defaultformat), "dateto": endTime.Format(defaultformat)}})
	return
}
func create(w http.ResponseWriter, r *http.Request) {
	randomsalt := RandStringBytesMaskImprSrc(7)
	h := encrypt_password("williamfl", randomsalt)
	insert := "INSERT INTO USERS (USERNAME,PASSWORD,SALT)"
	insert += "Values('williamfl',$1,$2)"
	req, err := SQL_exec(insert, []interface{}{hex.EncodeToString(h), randomsalt})
	log.Println(req)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprint(w, "Created")
}

func LoginEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		Login := request.(login)
		req, _ := SQL_exec("SELECT * FROM USERS WHERE username = $1", []interface{}{Login.username})
		if len(req) != 1 {
			return nil, GenError{"User not Found", http.StatusUnauthorized}
		}
		s := req[0]["salt"].(string)
		t := req[0]["password"].(string)
		check := USERS{s, t}
		if hex.EncodeToString(encrypt_password(Login.password, check.salt)) == check.password {
			//usersAllowed = append(usersAllowed, AuthorizeInfo{username: Login.username, authorizationId: []string{"Testing"}, authorizedUntil: time.Now().Add(15 * time.Minute)})
			return ResponseWithCookie{"Success", Login.username}, nil
		}
		return nil, GenError{"Bad Password", http.StatusUnauthorized}
	}
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Method != "POST" {
		return nil, GenError{"Method Not Allowed", http.StatusMethodNotAllowed}
	}
	username := r.FormValue("user")
	if username == "" {
		return nil, errors.New("Username blank")
	}
	password := r.FormValue("password")
	if password == "" {
		return nil, errors.New("Password blank")
	}
	return login{username, password}, nil
}

func UpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		if Last_update.IsZero() {
			Last_update = time.Now()
			// Now, we take the delay, and the person's name, and make a WorkRequest out of them.
			work := WorkRequest{Name: "Update", Arguments: []string{"node", "server.js"}}

			// Push the work onto the queue.
			WorkQueue <- work
			return Response{"Updating"}, nil
		}
		duration := time.Since(Last_update)
		if duration.Minutes() >= 120 {
			work := WorkRequest{Name: "Update", Arguments: []string{"node", "server.js"}}
			Last_update = time.Now()
			// Push the work onto the queue.
			WorkQueue <- work
			return Response{"Updating"}, nil
		}
		return Response{fmt.Sprintf("Cannot Update please wait %v Minutes", math.Floor((5-duration.Minutes())+.5))}, nil
	}
}
func RetrieveEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		v, err := svc.Retrieve(request.(RetrieveParameters))
		if err != nil {
			return nil, err
		}
		return RetrieveRequest{v}, nil
	}
}
func decodeRetrieveRequest(c context.Context, r *http.Request) (interface{}, error) {
	request := r.URL.Query()
	bad_parameters := make([]string, 0)
	cookieUser, err := r.Cookie("session")
	if err != nil {
		return nil, GenError{"No Current Session is good aaa", http.StatusUnauthorized}
	}
	if _, ok := request["dateto"]; !ok {
		bad_parameters = append(bad_parameters, "dateto")
	}
	if _, ok := request["datefrom"]; !ok {
		bad_parameters = append(bad_parameters, "datefrom")
	}
	if len(bad_parameters) > 0 {
		return nil, errors.New("Missing one or more parameters: " + strings.Join(bad_parameters, ", "))
	}
	return pass{cookieUser, RetrieveParameters{request.Get("dateto"), request.Get("datefrom")}}, nil
}
func Authorize() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			var temp AuthorizeInfo
			var ok bool
			//cookie_use, err := request.Cookie("session")
			/*if err != nil {
				return nil, GenError{"No Current Session is good", http.StatusUnauthorized}
			}*/
			var re = request.(pass)
			values := make(map[string]string)
			err := cookieHandler.Decode("session", re.cookie.Value, &values)
			if err != nil {
				log.Println("WHATTTTTTTTT")
			}
			//log.Println(ctx.Value("cookie").(http.Cookie).Value)
			if temp, ok = usersAllowed[values["name"]]; !ok {
				return nil, GenError{"No Current Session is good", http.StatusUnauthorized}
			}
			if _, ok = temp.authorizationId[values["session"]]; !ok {
				return nil, GenError{"No Current Session is good", http.StatusUnauthorized}
			}
			if !time.Now().Before(temp.authorizedUntil) {
				delete(usersAllowed, values["name"])
				return nil, GenError{"Session expired", http.StatusUnauthorized}
			}

			return next(ctx, re.Retrieve)
		}
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	//w.Header().set("Access-Control-Allow-Origin", "*")
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(response)
}
func encodeResponseWithCookie(_ context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(ResponseWithCookie)
	http.SetCookie(w, setSession(r.username))
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(r.Response)
}
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if strings.HasPrefix(err.Error(), "Missing one") {
		w.WriteHeader(http.StatusBadRequest)
	} else {

		switch erro := err.(type) {
		case GenError:
			w.WriteHeader(erro.getHttpError())
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Error", "error": err.Error(),
	})
}
