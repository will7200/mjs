package commands

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/will7200/mjs/api"
	"github.com/will7200/mjs/job"
)

var (
	disableLiveReload bool
	renderToDisk      bool
	serverAppend      bool
	serverInterface   string
	port              int
	serverWatch       bool
	verbose           bool
	Dispatch          *job.Dispatcher
	db                *gorm.DB
	err               error
)
var servercmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"run"},
	Short:   "Whip up a instance",
	RunE:    server,
}

func init() {
	servercmd.Flags().IntVarP(&port, "port", "p", 4004, "port on which to listen to")
	servercmd.Flags().BoolVar(&verbose, "verbose", false, "output log verbose")
	servercmd.Flags().String("dbname", "sqlite3", "database type")
	servercmd.Flags().String("connection", "./temp_db.db", "database connection string")
	servercmd.Flags().Int("workers", 4, "amount of workers in pool")
	viper.BindPFlag("port", servercmd.Flags().Lookup("port"))
	viper.BindPFlag("database.dbname", servercmd.Flags().Lookup("dbname"))
	viper.BindPFlag("database.connection", servercmd.Flags().Lookup("connection"))
	viper.BindPFlag("interface.workers", servercmd.Flags().Lookup("workers"))
}
func server(cmd *cobra.Command, args []string) error {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	var parsedPort string
	if port != 0 {
		parsedPort = fmt.Sprintf(":%d", port)
	} else {
		parsedPort = ":4004"
	}
	Dispatch = &job.Dispatcher{}
	Dispatch.StartDispatcher(viper.GetInt("interface.workers"))
	db, err = gorm.Open(viper.GetString("database.dbname"), viper.GetString("database.connection"))
	if err != nil {
		panic(fmt.Sprintf("failed to connect database \ntype %s with connection %s", viper.GetString("dbname"), viper.GetString("connection")))
	}
	db.LogMode(true)
	Dispatch.SetPersistStorage(db)
	fmt.Println(db.AutoMigrate(&job.Job{}, &job.JobStats{}).GetErrors())
	log.Infof("Starting Server on port %d", port)
	return api.StartServer(parsedPort, Dispatch, db)
}
