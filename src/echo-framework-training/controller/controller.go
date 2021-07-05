package controller

import (
	"echo-framework-training/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

// HandleIndexGet は Index のGet時のHTMLデータ生成処理を行います。
func HandleIndexGet(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

// 疎通確認用 : 後で消す
func HandeAPITestGet(c echo.Context) error {
	fmt.Println("server conection comming.......")
	return c.JSON(http.StatusOK, map[string]interface{}{"HttpResponse": "OK"})
}

func HandleAPIMasterDataGet(c echo.Context) error {
	var master = model.GetMasterData()
	return c.JSON(http.StatusOK, master)
}

func HandleAPIUserLogInPOST(c echo.Context) error {
	param := new(model.LoginData)
	var text string
	if err := c.Request().ParseForm(); err != nil {
		println("Formがとれないエラー")
	}
	for k, v := range c.Request().Form {
		fmt.Printf("%v : %v\n", k, v)
		text = k
	}
	text = strings.Replace(text, "{\"json\":", "", -1)
	text = strings.Replace(text, "} : []", "", -1)
	text = strings.Replace(text, "}}", "}", -1)
	println(text)
	var map1 map[string]string
	err := json.Unmarshal([]byte(text), &map1)
	if err != nil {
		fmt.Println(err)
	}

	if err := c.Bind(param); err != nil {
		return err
	}
	param.Email = map1["mail"]
	param.PassWord = map1["pass"]
	user, result := model.UserLogIn(param)
	return c.JSON(http.StatusOK, map[string]interface{}{"result": result, "loginData": user})
}

func HandleAPIStockListGet(c echo.Context) error {
	var json = model.GetJsonFromStockTable()
	return c.JSON(http.StatusOK, json)
}

func HandleAPIResisterPOST(c echo.Context) error {
	param := new(model.ResisterInfo)
	var text string
	if err := c.Request().ParseForm(); err != nil {
		println("Formがとれないエラー")
	}
	for k, v := range c.Request().Form {
		fmt.Printf("%v : %v\n", k, v)
		text = k
	}
	text = strings.Replace(text, "{\"json\":", "", -1)
	text = strings.Replace(text, "; []", "", -1)
	text = strings.Replace(text, "}}", "}", -1)
	println(text)
	var map1 map[string]string
	err := json.Unmarshal([]byte(text), &map1)
	if err != nil {
		fmt.Println(err)
	}

	println("member name is " + map1["name"])

	param.UserId = map1["user_id"]
	param.Member = map1["name"]
	param.Venue = map1["venue"]
	param.Type = map1["type"]
	param.Department = map1["department"]
	param.Amount = map1["amount"]

	isSuccess := model.ResisterRecord(param)
	return c.JSON(http.StatusOK, map[string]interface{}{"result": isSuccess})
}

func HandleAPIReservedListGet(c echo.Context) error {
	// TODO userIdを取得する
	userId := "100000"
	list := model.GetReservedList(userId)
	return c.JSON(http.StatusOK, list)
}

func HandleAPIRemoveRecordPOST(c echo.Context) error {
	// userIdを取得する
	userId := 100000
	var text string
	if err := c.Request().ParseForm(); err != nil {
		println("Formがとれないエラー")
	}
	for k, v := range c.Request().Form {
		fmt.Printf("%v : %v\n", k, v)
		text = k
	}
	text = strings.Replace(text, "{\"json\":", "", -1)
	text = strings.Replace(text, " : []", "", -1)
	text = strings.Replace(text, "}}", "}", -1)
	println(text)
	var map1 map[string]string
	err := json.Unmarshal([]byte(text), &map1)
	if err != nil {
		fmt.Println(err)
	}

	reserveId := map1["reserve_id"]
	recordId := map1["record_id"]
	println(reserveId)
	println(recordId)
	result := model.RemoveReservedRocord(reserveId, recordId, userId)
	return c.JSON(http.StatusOK, map[string]interface{}{"result": result})
}

func HandleAPIReservePOST(c echo.Context) error {
	param := new(model.ReserveInfo)
	var text string
	if err := c.Request().ParseForm(); err != nil {
		println("Formがとれないエラー")
	}
	for k, v := range c.Request().Form {
		fmt.Printf("%v : %v\n", k, v)
		text = k
	}

	text = strings.Replace(text, "{\"json\":", "", -1)
	text = strings.Replace(text, "}}", "}", -1)
	println(text)
	var map1 map[string]string
	err := json.Unmarshal([]byte(text), &map1)
	if err != nil {
		fmt.Println(err)
	}
	param.Member = map1["name"]
	param.Venue = map1["venue"]
	param.Type = map1["type"]
	param.Department = map1["department"]

	if err := c.Bind(param); err != nil {
		return err
	}
	result := model.ReserveTicket(param)
	println(result)
	return c.JSON(http.StatusOK, map[string]interface{}{"result": result})
}

type SearchResult struct {
	Amount string `json:"amount"`
}

func HandleAPISearchPOST(c echo.Context) error {
	param := new(model.SearchParams)
	var text string
	if err := c.Request().ParseForm(); err != nil {
		println("Formがとれないエラー")
	}
	for k, v := range c.Request().Form {
		fmt.Printf("%v : %v\n", k, v)
		text = k
	}
	text = strings.Replace(text, "{\"json\":", "", -1)
	text = strings.Replace(text, "} : []", "", -1)
	text = strings.Replace(text, "}}", "}", -1)
	println(text)
	var map1 map[string]string
	err := json.Unmarshal([]byte(text), &map1)
	if err != nil {
		fmt.Println(err)
	}

	println("member name is " + map1["name"])

	param.Name = map1["name"]
	param.Venue = map1["venue"]
	param.Type = map1["type"]
	param.Department = map1["department"]

	if err := c.Bind(param); err != nil {
		return err
	}

	amountStr := model.SearchStockInfo(param)
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return err
	}
	isSuccess := amount > 0

	return c.JSON(http.StatusOK, map[string]interface{}{"result": isSuccess})
}

func HandleAPIUserDataPOST(c echo.Context) error {
	param := new(model.UserIdHolder)
	var text string
	if err := c.Request().ParseForm(); err != nil {
		println("Formがとれないエラー")
	}
	for k, v := range c.Request().Form {
		fmt.Printf("%v : %v\n", k, v)
		text = k
	}
	text = strings.Replace(text, "{\"json\":", "", -1)
	text = strings.Replace(text, "} : []", "", -1)
	text = strings.Replace(text, "}}", "}", -1)
	println(text)
	var map1 map[string]string
	err := json.Unmarshal([]byte(text), &map1)
	if err != nil {
		fmt.Println(err)
	}
	param.UserId = map1["user_id"]
	if err := c.Bind(param); err != nil {
		return err
	}
	data := model.GetUserData(param)
	return c.JSON(http.StatusOK, data)
}

func HandleAPICreateUserPOST(c echo.Context) error {
	param := new(model.CreateUserInfo)
	var text string
	if err := c.Request().ParseForm(); err != nil {
		println("Formがとれないエラー")
	}
	for k, v := range c.Request().Form {
		fmt.Printf("%v : %v\n", k, v)
		text = k
	}
	text = strings.Replace(text, "{\"json\":", "", -1)
	text = strings.Replace(text, "} : []", "", -1)
	text = strings.Replace(text, "}}", "}", -1)
	println(text)
	var map1 map[string]string
	err := json.Unmarshal([]byte(text), &map1)
	if err != nil {
		fmt.Println(err)
	}
	param.Mail = map1["mail"]
	param.PassWord = map1["password"]
	param.UserNameFirst = map1["first_name"]
	param.UserNameLast = map1["last_name"]
	param.PostCode = map1["post_code"]
	param.Address = map1["address"]
	param.Phone = map1["phone"]
	if err := c.Bind(param); err != nil {
		return err
	}
	result := model.CreateUser(param)

	return c.JSON(http.StatusOK, result)
}
