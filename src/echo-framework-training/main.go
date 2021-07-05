package main

import (
	"echo-framework-training/controller"
	"html/template"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// "github.com/jinzhu/gorm"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
)

// レイアウト適用済のテンプレートを保存するmap
var templates map[string]*template.Template

// Template はHTMLテンプレートを利用するためのRenderer Interfaceです。
type Template struct {
}

// Render はHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込みます。
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return templates[name].ExecuteTemplate(w, "layout.html", data)
}

func main() {
	// Echoのインスタンスを生成
	e := echo.New()
	// テンプレートを利用するためのRendererの設定
	t := &Template{}
	e.Renderer = t

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORSの設定追加
	// 厳密にCORSを設定する（以下公式のサンプルコード）
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins: []string{"http://127.0.0.1", "http://192.168.1.26"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
	}))
	// ここからCorsの設定

	// 静的ファイルのパスを設定
	e.Static("/public/css/", "./public/css")
	e.Static("/public/js/", "./public/js/")
	e.Static("/public/img/", "./public/img/")

	// 各ルーティングに対するハンドラを設定
	e.GET("/", controller.HandleIndexGet)
	e.GET("/api/test", controller.HandeAPITestGet)
	e.GET("/api/stock", controller.HandleAPIStockListGet)
	e.POST("/api/search", controller.HandleAPISearchPOST)
	e.POST("/api/resister", controller.HandleAPIResisterPOST)
	e.POST("/api/reserve", controller.HandleAPIReservePOST)
	e.POST("/api/remove", controller.HandleAPIRemoveRecordPOST)
	e.GET("/api/reserved", controller.HandleAPIReservedListGet)
	e.POST("/api/login", controller.HandleAPIUserLogInPOST)
	e.POST("/api/createUser", controller.HandleAPICreateUserPOST)
	e.GET("/api/master_data", controller.HandleAPIMasterDataGet)
	e.POST("/api/user_data", controller.HandleAPIUserDataPOST)

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "hawksnowlog" && password == "secret" {
			return true, nil
		}
		return false, nil
	}))
	// basic auth
	track := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			println("request to /users")
			return next(c)
		}
	}
	g.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "/admin/users")
	}, track)

	// サーバーを開始
	e.Logger.Fatal(e.Start(":3000"))
}

// 初期化を行います。
func init() {
	loadTemplates()
}

// 各HTMLテンプレートに共通レイアウトを適用した結果を保存します（初期化時に実行）。
func loadTemplates() {
	var baseTemplate = "templates/layout.html"
	templates = make(map[string]*template.Template)
	templates["index"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/hello.html"))
	templates["hello_form"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/hello_form.html"))
}
