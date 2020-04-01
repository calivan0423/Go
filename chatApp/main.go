package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	// 웹 어플에서 세션을 생서 관리
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"

	//몽고db 패키지
	"gopkg.in/mgo.v2"
)

const (
	// session key information
	sessionKey    = "simple_chat_session"
	sessionSecret = "simple_chat_session_secret"
)

var (
	renderer     *render.Render
	mongoSession *mgo.Session
)

func init() {
	//create render
	renderer = render.New()

	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	mongoSession = s

}

func main() {
	//create router
	router := httprouter.New()

	//define handler

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "simple chat"})
	})

	router.GET("/rooms", createRoom)
	router.POST("/rooms", retrieveRooms)
	router.GET("/rooms/:id/messages", retrieveMessage)

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//로그인 페이지 렌더링
		renderer.HTML(w, http.StatusOK, "login", nil)
	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//세선에서 사용자 정보 제거 후 로그인 페이지 이동
		sessions.GetSession(r).Delete(currentUserKey)
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	router.GET("/auth/:acion/:provider", LoginHandler)

	//create negroni middleware
	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))

	//enroll router as handlerer to negroni
	n.UseHandler(router)

	//execute web server
	n.Run(":3000")

}
