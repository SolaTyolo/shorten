package main

import (
	"errors"
	"go_server/config"
	"go_server/model"
	"go_server/router"
	"go_server/router/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "go_server cofig file path")
)

// create gin http server
func main() {

	pflag.Parse()

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// set gin server mode and start server engine
	// debug ,release, test
	gin.SetMode(viper.GetString("runmode"))
	g := gin.New()
	router.Load(
		g,

		middleware.RequestId(),
	)

	// auto ping self server  to make sure router is working
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("the router has no reponse")
		}
		log.Println("the Router has deploy success")
	}()

	log.Printf("start to linstening the requests on the http address: %s", viper.GetString("addr"))
	log.Println(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Can't connect to the router")
}
