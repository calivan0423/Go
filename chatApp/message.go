package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"
)

const messageFetchSize = 10

type Message struct {
	ID        bson.ObjectId `bson: "_id" json: "id"`
	RoomId    bson.ObjectId `bson: "room_id" json: "room_id"`
	Content   string        `bson: "content" json: "content"`
	CreatedAt time.Time     `bson: "created_at" json: "created_at"`
	User      *User         `bson: "user" json:"user"`
}

func (m *Message) create() error {
	//몽고db 생성
	session := mongoSession.Copy()
	defer session.Close()

	//몽고db 아이디 생성
	m.ID = bson.NewObjectId()
	//메세지 생성기간 기록
	m.CreatedAt = time.Now()
	//message 정보 저장을 위한 몽고db 컬렌션 객체 생성
	c := session.DB("test").C("message")

	//message 컬렉션에 message 정보 저장
	if err := c.Insert(m); err != nil {
		return err
	}
	return nil

}

func retrieveMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = messageFetchSize
	}

	var messages []Message
	// _id 역순으로 정렬하여 limit 수만큼 message 조회
	err = session.DB("test").C("message").Find(bson.M{"room_id": bson.ObjectIdHex(ps.ByName("id"))}).Sort("_id").Limit(limit).All(&messages)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}

	renderer.JSON(w, http.StatusOK, messages)

}
