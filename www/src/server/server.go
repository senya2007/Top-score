package server

import (
	"log"
	"encoding/json"
	"net/http"
	"database/sql"
	"fmt"
	"time"
	"errors"
)

type NameScoreStruct struct {
	NameStruct
	Score string
}

type playerStruct struct {
	NamePasswordStruct
	Score string
}

type NamePasswordStruct struct{
	NameStruct
	Password string
}

type NameStruct struct{
	Name string
}

type TypeResponse int

const(
	ERROR TypeResponse = iota
	GOOD
	PLAYERS
)

func (e TypeResponse) String() string {
	switch e {
	case ERROR:
		return "error"
	case GOOD:
		return "good"
	case PLAYERS:
		return "players"
	}
	return ""
}
type ResponseJson struct {
	Type TypeResponse
	Value interface{}
	MethodName string
}

func GetJsonByteArrayFromPlayers(players []NameScoreStruct)([]byte,error) {
	if len(players) == 0{
		return nil,errors.New("Empty name struct")
	}
	result := struct {
		Type string
		Value []NameScoreStruct
		MethodName string
	}{Type:PLAYERS.String(), Value:players, MethodName:"updateAllScores"}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return  resultBytes, nil
}
func GetByteArrayFromResponse(responseJson ResponseJson) ([]byte,error){
	if responseJson == (ResponseJson{}){
		return nil, errors.New("Empty response json")
	}

	anonstr := struct {
		Type string
		Value string
		MethodName string
	}{Type:responseJson.Type.String(), Value:responseJson.Value.(string), MethodName:responseJson.MethodName}
	resultBytes, err := json.Marshal(anonstr)
	if err != nil {
		return nil,err
	}
	return resultBytes,nil
}


func ContainsPlayerInDB(player NamePasswordStruct) (int, error){
	if player == (NamePasswordStruct{}){
		return -1, errors.New("Empty player")
	}
	var query string
	query = "SELECT `id` FROM `players` " +
		"WHERE `name` = \""+player.Name+"\""
	result, err := db.Query(query)
	if err != nil {
		return -1, err
	}

	var id int
	for result.Next() {
		err = result.Scan(&id)

		if err != nil {
			return -1, err
		}
	}

	if id > 0{
		return  id,nil
	} else {
		return  -1,nil
	}
}

func ContainsNamePasswordInDB(player NamePasswordStruct) (int, error){
	var query string
	query = "SELECT `id` FROM `players` " +
		"WHERE `name` = \""+player.Name+"\" AND " +
		"`password` = \"" + player.Password+"\""
	result, err := db.Query(query)
	if err != nil {
		return -1, err
	}

	var id int
	for result.Next() {
		err = result.Scan(&id)

		if err != nil {
			return -1, err
		}
	}

	if id > 0{
		return  id,nil
	} else {
		return  -1,nil
	}
}

func Login(res http.ResponseWriter, req *http.Request)  {
	req.ParseForm()
	//LOG: map[{"test": "that"}:[]]
	var player NamePasswordStruct
	for key, _ := range req.Form {
		log.Println(key)
		//LOG: {"test": "that"}
		err := json.Unmarshal([]byte(key), &player)
		if err != nil {
			log.Println("Ошибка десериализации")
			result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:err.Error(), MethodName:"login"})
			res.Write(result)
			return
		}
	}

	var idPlayer int
	var err error

	idPlayer, err = ContainsNamePasswordInDB(player)

	if err != nil{
		log.Println("Ошибка проверки в базе пользователя ",player.Name)
		result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:err.Error(), MethodName:"login"})
		res.Write(result)
		return
	}

	var resultJsonBytes []byte
	if idPlayer > 0{
		log.Println("Успешно вошел пользователь ",player.Name)
		goodJson := ResponseJson{Value:"contains",Type:GOOD, MethodName:"login"}
		resultJsonBytes,err = GetByteArrayFromResponse(goodJson)
	} else {
		log.Println("Ошибка логина или пароля для ",player.Name)
		errorJson := ResponseJson{Value:"Ошибка логина или пароля ",Type:ERROR, MethodName:"login"}
		resultJsonBytes, err = GetByteArrayFromResponse(errorJson)
	}

	if err != nil{
		result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:err.Error(), MethodName:"login"})
		res.Write(result)
	}
	res.Write(resultJsonBytes)
}

var db *sql.DB

func SetDb(newDb *sql.DB){
	db = newDb
}
func CloseConnectDB()  {
	db.Close()
}
func ConnectToDB()(*sql.DB, error){
	db, err := sql.Open("mysql", "root@/topplayers")
	if err != nil {
		return nil ,err
	}
	return db, nil
}

func GetAllPlayers(db *sql.DB) ([]NameScoreStruct, error){
	if db == nil{
		return nil, errors.New("Db is null")
	}
	var players []NameScoreStruct
	res, err := db.Query("SELECT * FROM `players`")
	if err != nil {
		return nil, err
	}

	var id, name, password, score []byte

	for res.Next() {
		err = res.Scan(&id, &name, &password, &score)

		if err != nil {
			return nil, err
		}

		var temp NameScoreStruct
		temp.Name= string(name)
		temp.Score = string(score)
		players = append(players, temp)

		// Use the string value
		fmt.Println(string(name), string(score))
	}
	return players, nil
}

var players []NameScoreStruct

func UpdateScore(res http.ResponseWriter, req *http.Request){
	if req.Method == http.MethodPost {
		req.ParseForm()
		//LOG: map[{"test": "that"}:[]]
		var player playerStruct
		for key, _ := range req.Form {
			log.Println(key)
			//LOG: {"test": "that"}
			err := json.Unmarshal([]byte(key), &player)
			if err != nil {
				log.Println("Ошибка десериализации")
				result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка десериализации", MethodName:"updateScore"})
				res.Write(result)
				return
			}
		}

		var idPlayer int
		var errorInContainsPlayerInDb error
		var playerForCheckInDb NamePasswordStruct

		playerForCheckInDb.Name = player.Name
		playerForCheckInDb.Password = player.Password

		idPlayer, errorInContainsPlayerInDb = ContainsNamePasswordInDB(playerForCheckInDb)
		if errorInContainsPlayerInDb != nil{
			log.Println("Ошибка проверки в базе пользователя ", playerForCheckInDb.Name)
			result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка проверки пользователя в базе", MethodName:"updateScore"})
			res.Write(result)
			return
		}

		if idPlayer > 0{
			var err error
			_, err = db.Exec(
				"UPDATE `players` SET `score` = ? WHERE `id` = ?",
				player.Score,
				idPlayer,
			)
			time.Sleep(500 * time.Millisecond)
			if err != nil {
				log.Println("Ошибка обновления в базе")
				result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка обновления в базе", MethodName:"updateScore"})
				res.Write(result)
				return
			}else {
				log.Println("Обновлены очки для плеера ",player.Name)
				result,_ := GetByteArrayFromResponse(ResponseJson{Type:GOOD, Value:"true", MethodName:"updateScore"})
				res.Write(result)
				return
			}
		}
	}
}

func CreateUsers(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		req.ParseForm()
		//LOG: map[{"test": "that"}:[]]
		var player NamePasswordStruct
		for key, _ := range req.Form {
			log.Println(key)
			//LOG: {"test": "that"}
			err := json.Unmarshal([]byte(key), &player)
			if err != nil {
				log.Println("Ошибка десериализации")
				result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка десериализации", MethodName:"create"})
				res.Write(result)
				return
			}
		}

		var idPlayer int
		var errorInContainsPlayerInDb error

		idPlayer, errorInContainsPlayerInDb = ContainsPlayerInDB(player)
		if errorInContainsPlayerInDb != nil{
			log.Println("Ошибка проверки пользователя в базе")
			result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка проверки пользователя в базе", MethodName:"create"})
			res.Write(result)
			return
		}

		if idPlayer > 0{
			log.Println("player already created")
			result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Такой пользователь уже существует", MethodName:"create"})
			res.Write(result)
			return
		} else {
			var err error
			_, err = db.Exec(
				"INSERT INTO players (name, password) VALUES (?, ?)",
				player.Name,
				player.Password,
			)
			if err != nil {
				log.Println("Ошибка записи пользователя в базу")
				result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка записи пользователя в базу", MethodName:"create"})
				res.Write(result)
				return
			}else {
				log.Println("create new player")
				result,_ := GetByteArrayFromResponse(ResponseJson{Type:GOOD, Value:"create new player", MethodName:"create"})
				res.Write(result)
				return
			}
		}
	}
}

func UpdateAllScores(res http.ResponseWriter, req *http.Request){
	if req.Method == http.MethodPost {
		log.Println("Пользователи:")
		var errorFromGetAllPlayers error
		players, errorFromGetAllPlayers := GetAllPlayers(db)
		if(errorFromGetAllPlayers != nil){
			log.Println("Ошибка получения пользователей")
			result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка получения пользователей", MethodName:"updateAllScores"})
			res.Write(result)
			return
		}

		resultJsonBytes, err := GetJsonByteArrayFromPlayers(players)

		if err != nil{
			log.Println("Ошибка преобразования в JSON")
			result,_ := GetByteArrayFromResponse(ResponseJson{Type:ERROR, Value:"Ошибка преобразования в JSON", MethodName:"updateAllScores"})
			res.Write(result)
			return
		}
		res.Write(resultJsonBytes)
	}
}

