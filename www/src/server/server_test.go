package server

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"net/http/httptest"
)

func TestGetJsonByteArrayFromPlayers(t *testing.T){
	var players []NameScoreStruct

	var temp1 NameScoreStruct
	temp1.Name = "Bob"
	temp1.Score = "1000"
	var temp2 NameScoreStruct
	temp2.Name = "Jimmy"
	temp2.Score = "960"
	var temp3 NameScoreStruct
	temp3.Name = "Tom"
	temp3.Score = "900"
	players = append(players, temp1)
	players = append(players, temp2)
	players = append(players, temp3)

	byteArray,_ := GetJsonByteArrayFromPlayers(players)
	if(len(byteArray)==0 || byteArray == nil){
		t.Error("Players not parsing")
	}

	emptyByteArray, err := GetJsonByteArrayFromPlayers([]NameScoreStruct{})
	if(err.Error() != "Empty name struct" && emptyByteArray != nil){
		t.Error("Error parse empty players")
	}
}

func TestGetByteArrayFromResponse(t *testing.T)  {
	//byte,err := GetByteArrayFromResponse(ResponseJson{Type:GOOD, Value:string("test method"),MethodName:"testMethod" })

	byte,err := GetByteArrayFromResponse(ResponseJson{})
	if(err.Error()!="Empty response json" && byte != nil){
		t.Error("Error parsing")
	}
	}

func TestContainsPlayerInDB(t *testing.T)  {
	db,errToConnect := ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}

	SetDb(db)

	defer CloseConnectDB()

	playerId,err:= ContainsPlayerInDB(NamePasswordStruct{})

	if(playerId != -1 && err.Error() != "Empty player"){
		t.Error("Error get player")
	}
}

func TestContainsNamePasswordInDB(t *testing.T)  {
	db,errToConnect := ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}

	SetDb(db)

	defer CloseConnectDB()

	result, err:=ContainsNamePasswordInDB(NamePasswordStruct{})

	if(result != -1 && err != nil){
		t.Error("Error get id from DB")
	}
}

func TestGetAllPlayers(t *testing.T){
	db,errToConnect := ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}

	SetDb(db)

	defer CloseConnectDB()

	players, err:=GetAllPlayers(nil)
	if players != nil && err.Error() != "Db is null"{
		t.Error("Error parse from Db")
	}
}

func TestLogin(t *testing.T){
	db,errToConnect := ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}

	SetDb(db)

	defer CloseConnectDB()

	req := httptest.NewRequest("POST", "http://localhost:8080/login", nil)
	w := httptest.NewRecorder()
	Login(w,req)
}

func TestUpdateScore(t *testing.T){
	db,errToConnect := ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}

	SetDb(db)

	defer CloseConnectDB()

	req := httptest.NewRequest("POST", "http://localhost:8080/login", nil)
	w := httptest.NewRecorder()
	UpdateScore(w,req)
}

func TestCreateUsers(t *testing.T){
	db,errToConnect := ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}

	SetDb(db)

	defer CloseConnectDB()

	req := httptest.NewRequest("POST", "http://localhost:8080/login", nil)
	w := httptest.NewRecorder()
	CreateUsers(w,req)
}

func TestUpdateAllScores(t *testing.T){
	db,errToConnect := ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}

	SetDb(db)

	defer CloseConnectDB()

	req := httptest.NewRequest("POST", "http://localhost:8080/login", nil)
	w := httptest.NewRecorder()
	UpdateAllScores(w,req)
}