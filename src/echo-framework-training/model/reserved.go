package model

type ReservedRecordRow struct {
	RecordID     string `json:"record_id"`
	UserRecordID string `json:"user_record_id"`
	UserID       string `json:"user_id"`
	UserCode     string `json:"user_code"`
	VenueID      string `json:"venue"`
	MemberID     string `json:"member"`
	TypeID       string `json:"type"`
	DepartmentID string `json:"department"`
}

type RecordInfo struct {
	RecordID     string `json:"record_id"`
	UserRecordID string `json:"user_record_id"`
	UserID       string `json:"user_id"`
	UserCode     string `json:"user_code"`
	Venue        string `json:"venue"`
	Member       string `json:"member"`
	Type         string `json:"type"`
	Department   string `json:"department"`
}

// 予約一覧を取得するAPI
func GetReservedList(UserID string) interface{} {
	db := ConnectDB()
	rows, err := db.Query("select * from db_server.reserve_table where user_id=" + UserID)
	if err != nil {
		panic(err.Error())
	}
	master := GetMasterData()
	var records []RecordInfo
	for rows.Next() {
		row := ReservedRecordRow{}
		info := RecordInfo{}
		if err := rows.Scan(&row.RecordID, &row.UserRecordID, &row.UserID, &row.UserCode, &row.MemberID, &row.VenueID, &row.TypeID, &row.DepartmentID); err != nil {
			panic(err.Error())
		}
		info.RecordID = row.RecordID
		info.UserRecordID = row.UserRecordID
		info.UserID = row.UserID
		info.UserCode = row.UserCode
		for i := 0; i < len(master.MemberMaster); i++ {
			if row.MemberID == master.MemberMaster[i].MemberId {
				info.Member = master.MemberMaster[i].Name
			}
		}

		for i := 0; i < len(master.VenueMaster); i++ {
			if row.VenueID == master.VenueMaster[i].VenueId {
				info.Venue = master.VenueMaster[i].Venue_STR
			}
		}

		for i := 0; i < len(master.TypeMaster); i++ {
			if row.TypeID == master.TypeMaster[i].TypeId {
				info.Type = master.TypeMaster[i].Type_STR
			}
		}

		for i := 0; i < len(master.DepartmentMaster); i++ {
			if row.DepartmentID == master.DepartmentMaster[i].DepartmentId {
				info.Department = master.DepartmentMaster[i].Department_STR
			}
		}
		records = append(records, info)
	}
	return records
}
