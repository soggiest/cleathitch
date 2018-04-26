package main

import (
	"flag"
	"net/http"

	"github.com/soggiest/cleathitch/connector"

	"fmt"

	jwt "github.com/dgrijalva/jwt-go"

	k8sauth "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/soggiest/cleathitch/config"
)

const (
	v1Prefix = "k8s-aws-v1."
)

var (
	cfg config.Config
)

type Input struct {
	Apiversion string `json:"apiVersion" form:"apiVersion" query:"apiVersion"`
	Kind       string `json:"kind" form:"kind" query:"kind"`
	Spec       struct {
		Token string `json:"token" form:"token" query:"token"`
	} `json:"spec" form:"spec" query:"spec"`
}

func homeHandler(c echo.Context) (err error) {

	inputCH := new(Input)

	if err = c.Bind(inputCH); err != nil {
		//e.Logger(err)
		return
	}

	jwtToken, err := parseToken(inputCH.Spec.Token)
	if err != nil {
		fmt.Printf("ERROR-PARSE: %v\n", err)
	}

	jwtClaims := jwtToken.Claims.(jwt.MapClaims)

	//TODO: Make this part of the config, it should be defineable by the user
	username := fmt.Sprint(jwtClaims["name"])
	groups := connector.GetGroups(cfg, username)

	k8sauth := k8sauth.TokenReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: k8sauth.SchemeGroupVersion.String(),
			Kind:       "TokenReview",
		},
		Status: k8sauth.TokenReviewStatus{
			Authenticated: true,
			User: k8sauth.UserInfo{
				Username: username,
				UID:      "42",
				Groups:   groups,
			},
		},
	}
	return c.JSON(http.StatusOK, k8sauth)
}

func parseToken(idToken string) (*jwt.Token, error) {
	token, _ := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("kid"), nil
	})
	return token, nil
}

func main() {

	cfgFile := flag.String("config", "/etc/cleathitch/config.yaml", "The config file to use.")
	flag.Parse()

	cfg = config.ReadConfig(*cfgFile)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {}))

	//e.File("/", "./index.html")
	e.POST("/", homeHandler)

	e.Logger.Fatal(e.StartTLS(":8086", "/etc/cleathitch-tls/ch-tls.crt", "/etc/cleathitch-tls/ch-tls.key"))
}
