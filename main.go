package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//"strings"
	//"os"
	_ "github.com/lib/pq"
)

const (
	host     = "baasu.db.elephantsql.com"
	port     = 5432
	user     = "gxtkbltv"
	password = "uxK0VpqjTUpfo5_ce1ihCo631Ld6Y7_U"
	dbname   = "gxtkbltv"
)

// TWLMSecrets handles the secrets from TAP
type TWLMSecrets struct {
	TWlmDbUsername string `yaml:"twlm_db_username"`
	TWlmDbPassword string `yaml:"twlm_db_password"`
	TWlmDbDatabase string `yaml:"twlm_db_database"`
	TWlmDbHost     string `yaml:"twlm_db_host"`
}

type getResponse struct {
	PostgresStatuses PostgresStatuses `json:"company_table"`
	Error            string           `json:"error_message"`
}

// PostgresStatus is status of PG tables
type PostgresStatus struct {
	Name    string `json:"company_name"`
	Age     string `json:"age"`
	Address string `json:"address"`
}

// PostgresStatuses is Array of PostgresStatus
type PostgresStatuses []PostgresStatus

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	})

	http.HandleFunc("/pgstatus", func(w http.ResponseWriter, r *http.Request) {

		//fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s",
			host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}

		//postgresStatus := PostgresStatus{TableName: relname,AutoVacum: last_autovacuum, AutoAnalyze: last_autoanalyze}
		postgresStatuses := []PostgresStatus{}
		iterator := 0
		sqlStatement := `select name, age, address from COMPANY;`
		rows, err := db.Query(sqlStatement)
		//defer rows.Close()
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return
		case nil:
			for rows.Next() {
				var name string
				var age string     // just for intialization
				var address string // just for intialization
				err = rows.Scan(&name, &age, &address)
				if err != nil {
					// handle this error
					panic(err)
				}
				postgresStatus := PostgresStatus{Name: name, Age: age, Address: address}
				//fmt.Println(PGStatus.relname, PGStatus.last_autovacuum, PGStatus.last_autoanalyze)
				//fmt.Fprintf(w, "Table Name"+PGStatus.relname, " ")
				postgresStatuses = append(postgresStatuses, postgresStatus)
				iterator++
			}
		default:
			panic(err)
		}

		//pgStatusJson, err := json.Marshal(PGStatus)
		//fmt.Println(pgStatusJson)
		defer db.Close()

		err = db.Ping()
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := getResponse{
			PostgresStatuses: postgresStatuses,
			Error:            "none",
		}
		pgResponseJSON, err := json.Marshal(response)
		w.Write(pgResponseJSON)
		//fmt.Println(html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8082", nil))

}
