package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//	"gopkg.in/ldap.v2"

	"github.com/labstack/echo"
)

func homeHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	return c.JSONPretty(http.StatusOK, body, "  ")
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

	//e.File("/", "./index.html")
	e.GET("/", homeHandler)

	e.Logger.Fatal(e.Start(":8086"))
}
