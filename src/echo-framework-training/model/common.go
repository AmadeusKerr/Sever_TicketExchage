package model

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MasterData struct {
	MemberMaster     []Member     `json:"member_master"`
	VenueMaster      []Venue      `json:"venue_master"`
	TypeMaster       []Type       `json:"type_master"`
	DepartmentMaster []Department `json:"department_master"`
}

type Member struct {
	ID       int    `json:"id"`
	MemberId string `json:"member_id"`
	Name     string `json:"name"`
}

type Venue struct {
	ID          int
	VenueId     string
	Date        string
	Location    string
	Venue_STR   string
	TalkLimit   int
	InstaxLimit int
	SignLimit   int
}

type Type struct {
	ID       int
	TypeId   string
	Type_STR string
}

type Department struct {
	ID             int
	DepartmentId   string
	Department_STR string
}

func ConnectDB() (db *sql.DB) {
	DBMS := "mysql"
	USER := "admin"
	PASS := "yamanoku"
	PROTOCOL := "tcp(database-1.cz2pgcfvn1mm.ap-northeast-1.rds.amazonaws.com:3306)"
	DBNAME := "db_server"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME

	db, err := sql.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("DB接続成功")
	}
	return db
}

func GetMasterData() MasterData {
	var master MasterData
	db := ConnectDB()
	defer db.Close()
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
	master.MemberMaster = members

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
	master.VenueMaster = venues

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
	master.TypeMaster = typeList

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
	master.DepartmentMaster = departments
	defer db.Close()
	return master
}

func ConvertSQLString(s string) string {
	str := "\"" + s + "\""
	return str
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	s = "\"" + s + "\""
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

type UserIdHolder struct {
	UserId string `json:"user_id"`
}

type UserData struct {
	UserID     string         `json:"user_id"`
	UserCode   string         `json:"user_code"`
	FirstName  sql.NullString `json:"first_name"`
	LastName   sql.NullString `json:"last_name"`
	PassWord   string         `json:"pass"`
	Mail       string         `json:"mail"`
	PostCode   sql.NullString `json:"post_code"`
	Address    sql.NullString `json:"address"`
	Phone      sql.NullString `json:"phone"`
	TalkPoint  int            `json:"talk_pt"`
	ChekiPoint int            `json:"cheki_pt"`
	SignPoint  int            `json:"sign_pt"`
	CreateAt   string         `json:"create_at"`
	UpdateAt   string         `json:"update_at"`
}

func GetUserData(holder *UserIdHolder) interface{} {
	db := ConnectDB()
	row := db.QueryRow("select * db_server.from usr_table where user_id=" + holder.UserId)
	data := UserData{}
	err := row.Scan(&data.UserID, &data.UserCode, &data.FirstName, &data.LastName, &data.PassWord,
		&data.Mail, &data.PostCode, &data.Address, &data.Phone, &data.TalkPoint, &data.ChekiPoint, &data.SignPoint, &data.CreateAt, &data.UpdateAt)
	if err != nil {
		println(err.Error())
	}
	return data
}

func GetDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	return now.Format(layout)
}
