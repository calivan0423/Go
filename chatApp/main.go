package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	// 웹 어플에서 세션을 생서 관리
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"

	//몽고db 패키지
	"gopkg.in/mgo.v2"

	//웹소켓 패키지
	"github.com/gorilla/websocket"
)

const (
	// session key information
	sessionKey       = "simple_chat_session"
	sessionSecret    = "simple_chat_session_secret"
	socketBufferSize = 1024
)

var (
	renderer     *render.Render
	mongoSession *mgo.Session
	//http커넥션이 웹소켓을 사용할 수 있도록
	upgrader = websocket.Upgrader{
		ReadBufferSize:  socketBufferSize,
		WriteBufferSize: socketBufferSize,
	}
)

func init() {
	//create render
	renderer = render.New()

	s, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)
	}
	//세선을 변화없는 동작으로 바꾸기위한 옵션
	s.SetMode(mgo.Monotonic, true)
	mongoSession = s

}

func main() {
	//create router
	router := httprouter.New()

	//define handler

	router.GET("/", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		renderer.HTML(w, http.StatusOK, "index", map[string]string{"title": "simple chat"})
	})

	router.GET("/rooms", retrieveRooms)
	router.GET("/rooms/:id", retrieveRoom)
	router.POST("/rooms", createRoom)
	router.DELETE("/rooms/:id", deleteRoom)

	router.GET("/rooms/:id/messages", retrieveMessages)

	router.GET("/info", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		u := GetCurrentUser(r)
		info := map[string]interface{}{"currrent_user": u, "clients": clients}
		renderer.JSON(w, http.StatusOK, info)
	})

	router.GET("/login", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//로그인 페이지 렌더링
		renderer.HTML(w, http.StatusOK, "login", nil)
	})

	router.GET("/logout", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//세선에서 사용자 정보 제거 후 로그인 페이지 이동
		sessions.GetSession(r).Delete(currentUserKey)
		http.Redirect(w, r, "/login", http.StatusFound)
	})

	router.GET("/auth/:action/:provider", loginHandler)

	router.GET("/ws/:room_id", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Printf("HELLO")
		conn, err := upgrader.Upgrade(w, r, nil)
		log.Printf("HELLO2")
		if err != nil {
			log.Fatal("ServeHTTP:", err)
			return
		}
		newClient(conn, ps.ByName("room_id"), GetCurrentUser(r))
	})

	//create negroni middleware
	n := negroni.Classic()
	store := cookiestore.New([]byte(sessionSecret))
	n.Use(sessions.Sessions(sessionKey, store))

	n.Use(LoginRequired("/login", "/auth"))

	//enroll router as handlerer to negroni
	n.UseHandler(router)

	//execute web server
	n.Run(":3000")

}
