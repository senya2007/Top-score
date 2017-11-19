package main

import (
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"server"
)

func main() {
	db,errToConnect := server.ConnectToDB()
	if(errToConnect!= nil){
		panic(errToConnect)
	}
	server.SetDb(db)

	defer server.CloseConnectDB()
	http.HandleFunc("/create", server.CreateUsers)
	http.HandleFunc("/updateScore", server.UpdateScore)
	http.HandleFunc("/updateAllScores", server.UpdateAllScores)
	http.HandleFunc("/login", server.Login)
	http.ListenAndServe(":8080", nil)
}