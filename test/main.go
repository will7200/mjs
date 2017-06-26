package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/will7200/mjs/api"
	"github.com/will7200/mjs/job"

	"github.com/codegangsta/cli"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

var (
	REPEAT   = "R2"
	INTERVAL = "PT1S"
)
var newjob = []byte(
	`{"Name":"Youtube Channel","ID":"11","Owner":"William"}`,
)
var db *gorm.DB
var err error

func main() {
	Dispatch := &job.Dispatcher{}
	app := cli.NewApp()
	app.Usage = "Microservice Job Service"
	app.Version = ".1"
	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "start",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port, p",
					Value: 4004,
					Usage: "Port for MJS to run on.",
				},
				cli.StringFlag{
					Name:  "config, c",
					Value: "",
					Usage: "Point to the Configuration File",
				},
				cli.StringFlag{
					Name:  "configPath, C",
					Value: "",
					Usage: "Path for the Configuration File",
				},
				cli.BoolFlag{
					Name:  "verbose, v",
					Usage: "Set for verbose logging.",
				},
			},
			Action: func(c *cli.Context) {
				if c.Bool("v") {
					log.SetLevel(log.DebugLevel)
				}
				var parsedPort string
				port := c.Int("port")
				if port != 0 {
					parsedPort = fmt.Sprintf(":%d", port)
				} else {
					parsedPort = ":4004"
				}
				Dispatch.StartDispatcher(4)
				db, err = gorm.Open("postgres", "postgres://williamfl:williamfl@192.168.1.25/wjob")
				if err != nil {
					panic("failed to connect database")
				}
				db.LogMode(true)
				Dispatch.SetPersistStorage(db)
				fmt.Println(db.AutoMigrate(&job.Job{}, &job.JobStats{}).GetErrors())
				log.Infof("Starting Server on port %d", port)
				log.Fatal(api.StartServer(parsedPort, Dispatch, db))
			},
		},
	}
	app.Run(os.Args)
}
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("~./.mjs/")
	viper.AddConfigPath("$HOME/.mjs")
	err := viper.ReadInConfig()
	if err != nil {
		//setDefaultsAll()
	}
}
