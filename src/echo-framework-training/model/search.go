package model

import (
	"database/sql"
	"strconv"
)

type SearchParams struct {
	Name       string `json:"name"`
	Venue      string `json:"venue"`
	Type       string `json:"type"`
	Department string `json:"department"`
}

type SearchRow struct {
	ID           int `json:"ID"`
	MemberId     int `json:"member_id"`
	VenueId      int `json:"venue_id"`
	TypeId       int `json:"type_id"`
	DepartmentId int `json:"department_id"`
	Amount       int `json:"amount"`
}

func SearchStockInfo(info *SearchParams) string {
	db := ConnectDB()
	master := GetMasterData()
	var member_id string
	for i := 0; i < len(master.MemberMaster); i++ {
		if info.Name == master.MemberMaster[i].Name {
			member_id = master.MemberMaster[i].MemberId
		}
	}
	var venue_id string
	for i := 0; i < len(master.VenueMaster); i++ {
		if info.Venue == master.VenueMaster[i].Venue_STR {
			venue_id = master.VenueMaster[i].VenueId
		}
	}
	var type_id string
	for i := 0; i < len(master.TypeMaster); i++ {
		if info.Type == master.TypeMaster[i].Type_STR {
			type_id = master.TypeMaster[i].TypeId
		}
	}
	var department_id string
	for i := 0; i < len(master.DepartmentMaster); i++ {
		if info.Department == master.DepartmentMaster[i].Department_STR {
			department_id = master.DepartmentMaster[i].DepartmentId
		}
	}
	text := "select * from db_server.stock_table where member_id=" + member_id + " && venue_id=" + venue_id +
		" && type_id=" + type_id + " && department_id=" + department_id
	println(text)
	result, err := db.Query(text)

	if err != nil {
		// 0件だった場合
		if err == sql.ErrNoRows {
			return "0"
		}
		// 0件の場合？？
		panic(err.Error())
	}
	var row []SearchRow
	for result.Next() {
		search := SearchRow{}
		if err := result.Scan(&search.ID, &search.MemberId, &search.VenueId, &search.TypeId, &search.DepartmentId, &search.Amount); err != nil {
			panic(err.Error())
		}
		row = append(row, search)
	}
	return strconv.Itoa(row[0].Amount)
}
