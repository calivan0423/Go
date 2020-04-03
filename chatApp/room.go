package main

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"
	//http 요청 내용을 구조체로 변환하기 위한 패키지
	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
)

type Room struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}

//request 데이터를 Room 타입 구조체로 변환하기 위하여 fieldMap 메서드 추가
func (r *Room) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{&r.Name: "name"}
}

//체팅방으로 정보를 생성 및 조회하는 REST API 작성
func createRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//binding 패키지로 room 생성 요청 정보를 Room 타입 값으로 변환
	r := new(Room)
	errs := binding.Bind(req, r)
	if errs.Handle(w) {
		return
	}

	//몽고db 세션 생성
	session := mongoSession.Copy()
	// 몽고db 세션을 단는 코드를 defer로 등록
	defer session.Close()

	//몽고db 아이디 생성
	r.ID = bson.NewObjectId()
	//room 정보 저장을 위한 몽고db 컬렉션 객체 생성
	c := session.DB("test").C("rooms")

	//rooms 컬렉션에 room 정보 저장
	if err := c.Insert(r); err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	//처리 결과 반환
	renderer.JSON(w, http.StatusCreated, r)

}

func retrieveRooms(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//몽고 db 세선 생성
	session := mongoSession.Copy()
	defer session.Close()

	var rooms []Room
	//모든 Room 저보 조회
	err := session.DB("test").C("rooms").Find(nil).All(&rooms)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, rooms)

}

func retrieveRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	var room Room
	err := session.DB("test").C("rooms").FindId(bson.ObjectIdHex(ps.ByName("id"))).One(&room)
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusOK, room)
}

func deleteRoom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	session := mongoSession.Copy()
	defer session.Close()

	err := session.DB("test").C("rooms").RemoveId(bson.ObjectIdHex(ps.ByName("id")))
	if err != nil {
		renderer.JSON(w, http.StatusInternalServerError, err)
		return
	}
	renderer.JSON(w, http.StatusNoContent, nil)
}
