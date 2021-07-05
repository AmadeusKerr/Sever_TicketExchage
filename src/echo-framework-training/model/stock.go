package model

import "strconv"

// 在庫情報テーブル情報
type StockInfos struct {
	Stock_Id      int    `json:"stock_id"`
	Member_Id     string `json:"name"`
	Venue_Id      string `json:"venue"`
	Type_Id       string `json:"type"`
	Department_Id string `json:"department"`
	Amount        string `json:"amount"`
}

type StockData struct {
	StockId    string `json:"stock_id"`
	Name       string `json:"name"`
	Venue      string `json:"venue"`
	Type       string `json:"type"`
	Department string `json:"department"`
	Amount     string `json:"amount"`
}

type StockDataRow struct {
	StockID      int `json:"id"`
	MemberID     int `json:"member"`
	VenueID      int `json:"venue"`
	TypeID       int `json:"type"`
	DepartmentID int `json:"department"`
	Amount       int `json:"amount"`
}

func GetJsonFromStockTable() interface{} {
	db := ConnectDB()
	m, err := db.Query("select * from db_server.member_master")
	if err != nil {
		panic(err.Error())
	}
	var members []Member
	for m.Next() {
		member := Member{}
		if err := m.Scan(&member.ID, &member.MemberId, &member.Name); err != nil {
			println(err)
		}
		members = append(members, member)
	}

	v, err := db.Query("select * from db_server.venue_master")
	if err != nil {
		panic(err.Error())
	}
	var venues []Venue
	for v.Next() {
		venue := Venue{}
		if err := v.Scan(&venue.ID, &venue.VenueId, &venue.Date, &venue.Location,
			&venue.Venue_STR, &venue.TalkLimit, &venue.InstaxLimit, &venue.SignLimit); err != nil {
			println(err)
		}
		venues = append(venues, venue)
	}

	t, err := db.Query("select * from db_server.type_master")
	if err != nil {
		panic(err.Error())
	}
	var typeList []Type
	for t.Next() {
		_type := Type{}
		if err := t.Scan(&_type.ID, &_type.TypeId, &_type.Type_STR); err != nil {
			println(err)
		}
		typeList = append(typeList, _type)
	}

	d, err := db.Query("select * from db_server.department_master")
	if err != nil {
		panic(err.Error())
	}
	var departments []Department
	for d.Next() {
		department := Department{}
		if err := d.Scan(&department.ID, &department.DepartmentId, &department.Department_STR); err != nil {
			println(err)
		}
		departments = append(departments, department)
	}

	json := map[string]interface{}{"State": "ok"}
	count := 1
	rows, err := db.Query("select * from db_server.stock_table")
	if err != nil {
		panic(err.Error())
	}
	var datas []StockData
	for rows.Next() {
		info := StockInfos{}
		data := StockData{}
		if err := rows.Scan(&info.Stock_Id, &info.Member_Id, &info.Venue_Id, &info.Type_Id, &info.Department_Id, &info.Amount); err != nil {
			println(err)
		}
		data.StockId = strconv.Itoa(info.Stock_Id)
		data.Amount = info.Amount + "枚"
		var venue string
		for i := 0; i < len(venues); i++ {
			if info.Venue_Id == venues[i].VenueId {
				venue = venues[i].Venue_STR
				json["venue"+strconv.Itoa(count)] = venue
				data.Venue = venue
				break
			}
		}

		var member string
		for i := 0; i < len(members); i++ {
			if info.Member_Id == members[i].MemberId {
				member = members[i].Name
				json["name"+strconv.Itoa(count)] = member
				data.Name = member
				break
			}
		}
		var types string
		for i := 0; i < len(typeList); i++ {
			if info.Type_Id == typeList[i].TypeId {
				types = typeList[i].Type_STR
				json["type"+strconv.Itoa(count)] = types
				data.Type = types
				break
			}
		}
		var department string
		for i := 0; i < len(departments); i++ {
			if info.Department_Id == departments[i].DepartmentId {
				department = departments[i].Department_STR
				json["department"+strconv.Itoa(count)] = department
				data.Department = department
				break
			}
		}
		datas = append(datas, data)
	}
	defer db.Close()
	return datas
}
