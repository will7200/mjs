package main

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Initilize() {
	var connInfo string = "postgres://williamfl:williamfl@192.168.1.25/infodb"
	var err error
	db, err = sql.Open("postgres", connInfo)
	db.SetMaxOpenConns(2)
	if err != nil {
		log.Print(err)
	}
}

// SQLExec will execute the string with the provide parameters
func SQLExec(cmd string, parameters []interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(cmd, parameters...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if cols == nil {
		return nil, errors.New("Emtpy")
	}
	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}
	var dataJ map[string]interface{}
	array := make([]map[string]interface{}, 0)
	for rows.Next() {
		dataJ = make(map[string]interface{})
		err = rows.Scan(vals...)
		if err != nil {
			log.Println(err)
		}
		for i := 0; i < len(vals); i++ {
			dataJ[cols[i]] = getValue(vals[i].(*interface{}))
		}
		array = append(array, dataJ)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return array, nil
}

func getValue(pval *interface{}) string {
	switch v := (*pval).(type) {
	case nil:
		return "NULL"
	case bool:
		if v {
			return "True"
		}
		return "False"
	case int64:
		return strconv.FormatInt(v, 10)
	case []byte:
		return string(v)
	case time.Time:
		return v.Format("2006-01-02")
	case string:
		return v
	default:
		return v.(string)
	}
}
