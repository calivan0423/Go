package main

import (
	"log"
	"net/http"
	"strings"

	sessions "github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/urfave/negroni"
)

const (
	nextPageKey     = "next_page" //세선에 저장되는 next page의 키
	authSecurityKey = "auth_security_key"
)

func init() {
	//gomniauth 정보 세팅
	gomniauth.SetSecurityKey(authSecurityKey)
	gomniauth.WithProviders(google.New("826357128682-7j3osck2fn7b9o1tl3bvk384shouvrsd.apps.googleusercontent.com", "Zd1sA9aPl0AFDrTW6zzH7RiL", "http://127.0.0.1:3000/auth/callback/google"))
	/*
		gomniauth.WithProviders(
			google.New("636296155193-a9abes4mc1p81752l116qkr9do6oev3f.apps.googleusercontent.com", "EVvuy0Agv4jWflml0pvC6-vI",
				"http://127.0.0.1:3000/auth/callback/google"),
		)
	*/
}

func LoginRequired(ignore ...string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		//ignore url 이면 다음 핸들러 실행
		for _, s := range ignore {
			if strings.HasPrefix(r.URL.Path, s) {
				next(w, r)
				return
			}
		}

		u := GetCurrentUser(r)

		//현재 유저 정보가 유효하면 만료시간을 갱신, 다음 핸들러 실행
		if u != nil && u.Valid() {
			SetCurrentUser(r, u)
			next(w, r)
			return
		}

		//현재 유저 정보가 유효하지 않으면 현재 유저를 nil로
		SetCurrentUser(r, nil)

		//로그인 후 이동할 url을 세선에 저장
		sessions.GetSession(r).Set(nextPageKey, r.URL.RequestURI())

		//로그인 페이지로 리다이렉트
		http.Redirect(w, r, "/login", http.StatusFound)

	}
}
func loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	action := ps.ByName("action")
	provider := ps.ByName("provider")
	s := sessions.GetSession(r)

	switch action {
	case "login":
		// gomniauth Provider의 login 페이지로 이동
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}
		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		http.Redirect(w, r, loginUrl, http.StatusFound)
	case "callback":
		// gomniauth 콜백 처리
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln(err)
		}

		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln(err)
		}

		// 콜백 결과로 부터 사용자 정보 확인
		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln(err)
		}

		if err != nil {
			log.Fatalln(err)
		}

		u := &User{
			Uid:       user.Data().Get("id").MustStr(),
			Name:      user.Name(),
			Email:     user.Email(),
			AvatarUrl: user.AvatarURL(),
		}

		SetCurrentUser(r, u) // 사용자 정보를 세션에 저장

		http.Redirect(w, r, s.Get(nextPageKey).(string), http.StatusFound)
	default:
		http.Error(w, "Auth action '"+action+"' is not supported", http.StatusNotFound)
	}
}
