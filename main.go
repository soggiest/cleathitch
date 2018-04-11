package main

import (
	"net/http"
	//	"gopkg.in/ldap.v2"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Input struct {
	Apiversion string `json:"apiVersion" form:"apiVersion" query:"apiVersion"`
	Kind       string `json:"kind" form:"kind" query:"kind"`
	Spec       struct {
		Token string `json:"token" form:"token" query:"token"`
	} `json:"spec" form:"spec" query:"spec"`
}

func homeHandler(c echo.Context) (err error) {
	//body, err := ioutil.ReadAll(c.Request().Body)
	inputCH := new(Input)

	if err = c.Bind(inputCH); err != nil {
		fmt.Printf("ERROR?! %v", err)
		//e.Logger(err)
		return
	}

	fmt.Printf("INPUT RECEIVED: %v\n", inputCH)

	return c.JSON(http.StatusOK, inputCH)
	//	fmt.Printf("BEARER TOKEN: %v\n", c.Body)

	/*getBearerToken := r.Body
	modifyBearerToken := bearertokenshit(getBearerToken)

	jsonBody := '{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind": "TokenReview",
		"status": {
		  "authenticated": true,
		  "user": {
			"username": "heptio",
			"uid": "42",
			"groups": [
			  "developers",
			  "qa"
			],
			"extra": {
			  "extrafield1": [
				"extravalue1",
				"extravalue2"
			  ]
			}
		  }
		}
	  }'

	  WRITETOHTTP(jsonBody, w)
	*/
}

func main() {
	//	cfgFile := flag.String("config", "", "The config file to use.")
	//	flag.Parse()

	//cfg := config.ReadConfig(*cfgFile)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {}))

	//e.File("/", "./index.html")
	e.POST("/", homeHandler)

	e.Logger.Fatal(e.StartTLS(":8086", "/etc/cleathitch/tls/ch-tls.crt", "/etc/cleathitch/tls/ch-tls.key"))
}
