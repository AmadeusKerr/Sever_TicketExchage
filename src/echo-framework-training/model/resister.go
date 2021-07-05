package model

import "strconv"

type ResisterInfo struct {
	UserId     string `json:"user_id"`
	Member     string `json:"member"`
	Venue      string `json:"venue"`
	Type       string `json:"type"`
	Department string `json:"department"`
	Amount     string `json:"amount"` // 不要かも
	Trade      string `json:"trade"`
}

// TODO 使われてない？
type ResisterRow struct {
	UserID     int `json:"user_id"`
	UserCode   int `json:"user_code"`
	Member     int `json:"member"`
	Venue      int `json:"venue"`
	Type       int `json:"type"`
	Department int `json:"department"`
	Amount     int `json:"amount"` // 不要かも
	Trade      int `json:"trade"`
}

type TypePoint struct {
	Point int `json:"have_point"`
}

// チケット登録処理
func ResisterRecord(info *ResisterInfo) bool {
	db := ConnectDB()
	master := GetMasterData()
	var member_id string
	for i := 0; i < len(master.MemberMaster); i++ {
		if info.Member == master.MemberMaster[i].Name {
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
	row, err := db.Query("select * from db_server.stock_table where member_id=" + member_id + " && venue_id=" + venue_id +
		" && type_id=" + type_id + " && department_id=" + department_id)
	if err != nil {
		println("First Query error")
		panic(err.Error())
	}
	// すでに該当レコードがあったら
	resister := StockDataRow{}
	isSuccess := false
	if !row.Next() {
		ins, err := db.Prepare("insert into db_server.stock_table value(user_id, user_code, member_id, venue_id, type_id, department_id, amount) VALUES(?, ?, ?, ?, ?)")
		if err != nil {
			println("Insert prepare error")
			panic(err.Error())
		}
		ins.Exec(info.Member, info.Venue, info.Type, info.Department, info.Amount)
		isSuccess = true
	} else {
		result := db.QueryRow("select * from db_server.stock_table where member_id=" + member_id + " && venue_id=" + venue_id + " && type_id=" + type_id + " && department_id=" + department_id)
		result.Scan(&resister.StockID, &resister.MemberID, &resister.VenueID, &resister.TypeID, &resister.DepartmentID, &resister.Amount)
		if err != nil {
			panic(err.Error())
		}
		amountNow := resister.Amount
		addAmount, err := strconv.Atoi(info.Amount)
		if err != nil {
			panic(err.Error())
		}
		amount := addAmount + amountNow
		resultAmount := strconv.Itoa(amount)

		ins, err := db.Prepare("update db_server.stock_table set amount=? where stock_id=?")
		if err != nil {
			panic(err.Error())
		}
		ins.Exec(resultAmount, strconv.Itoa(resister.StockID))
		isSuccess = true
		var point_type string
		switch resister.TypeID {
		case 1: //個別トーク会
			point_type = "talk_pt"
		case 2: // 個別チェキ会
			point_type = "cheki_pt"
		case 3: // 個別サイン会
			point_type = "sign_pt"
		}

		query := db.QueryRow("select " + point_type + " from db_server.usr_table where user_id=" + info.UserId)
		point := TypePoint{}
		err = query.Scan(&point.Point)
		if err != nil {
			println(err.Error())
		}
		additionalAmount, _ := strconv.Atoi(info.Amount)
		finallyPoint := point.Point + additionalAmount
		println(finallyPoint)
		upd, err := db.Prepare("update db_server.usr_table set " + point_type + "=" + strconv.Itoa(finallyPoint) + " where user_id=" + info.UserId)
		if err != nil {
			println(err.Error())
		}
		upd.Exec()
	}
	defer db.Close()
	return isSuccess
}
