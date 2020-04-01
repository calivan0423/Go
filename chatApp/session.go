package main

import (
	"encoding/json"
	"net/http"
	"time"

	sessions "github.com/goincremental/negroni-sessions"
)

const (
	currentUserKey  = "oauth2_current_user" // 현재 유저의 키
	sessionDuration = time.Hour             // 로그인 세션 유지기간

)

type User struct {
	Uid       string    `json:"uid"`
	Name      string    `json:"name"`
	Email     string    `json:"user"`
	AvatarUrl string    `json:"avata_rul"`
	Expired   time.Time `json:"expired"`
}

func (u *User) Valid() bool {
	//현재 시간 기준 만료시간
	return u.Expired.Sub(time.Now()) > 0
}

func (u *User) Refresh() {
	u.Expired = time.Now().Add(sessionDuration)
}

func GetCurrenUser(r *http.Request) *User {
	//세션어세 현재 유저 정보 get
	s := sessions.GetSession(r)

	if s.Get(currentUserKey) == nil {
		return nil
	}

	data := s.Get(currentUserKey).([]byte)
	var u User
	json.Unmarshal(data, &u) //첫번째 파라미터에는 JSON 데이타를, 두번째 파라미터에는 출력할 구조체(혹은 map)를 포인터로 지정한다
	return &u
}

func SetCurrentUser(r *http.Request, u *User) {
	if u != nil {
		//현재 유저 만료 시간 갱신
		u.Refresh()
	}

	//세션에 현재유저정보를 json으로 저장
	s := sessions.GetSession(r)
	val, _ := json.Marshal(u) //JSON으로 인코딩된 바이트배열과 에러객체를 리턴
	s.Set(currentUserKey, val)
}
