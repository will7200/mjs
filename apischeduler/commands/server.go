package commands

import (
	"fmt"
	stdlog "log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	colorable "github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/will7200/mjs/apischeduler/endpoints"
	apischedulerhttp "github.com/will7200/mjs/apischeduler/http"
	"github.com/will7200/mjs/apischeduler/service"
	"github.com/will7200/mjs/job"
	"github.com/will7200/mjs/utils/hooks"
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
	showHTTPDir       bool
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
	servercmd.Flags().String("homedir", "./mda/", "home directory to download into")
	servercmd.Flags().Int("workers", 4, "Amount of workers for the Dispatcher")
	servercmd.Flags().BoolVar(&showHTTPDir, "httpdir", false, "Output the http directory")
	viper.BindPFlag("interface.port", servercmd.Flags().Lookup("port"))
	viper.BindPFlag("database.dbname", servercmd.Flags().Lookup("dbname"))
	viper.BindPFlag("database.connection", servercmd.Flags().Lookup("connection"))
	viper.BindPFlag("interface.workers", servercmd.Flags().Lookup("workers"))
	viper.BindPFlag("interface.home", servercmd.Flags().Lookup("homedir"))
	viper.SetEnvPrefix("mda") // will be uppercased automatically
	viper.BindEnv("VERBOSE", "verbose")
}
func server(cmd *cobra.Command, args []string) error {
	verbose = viper.GetBool("verbose")
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	var parsedPort string
	port = viper.GetInt("interface.port")
	if port != 0 {
		parsedPort = fmt.Sprintf(":%d", port)
	} else {
		parsedPort = ":4004"
	}
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
	log.AddHook(hooks.NewHook(&hooks.CallerHookOptions{
		Field: "src",
		Flags: stdlog.Lshortfile,
	}))
	Dispatch = &job.Dispatcher{}
	Dispatch.StartDispatcher(viper.GetInt("interface.workers"))
	db, err = gorm.Open(viper.GetString("database.dbname"), viper.GetString("database.connection"))
	if verbose {
		db.LogMode(true)
	}
	if err != nil {
		panic(fmt.Sprintf("failed to connect database \ntype %s with connection %s", viper.GetString("database.dbname"), viper.GetString("database.connection")))
	}
	if errors := db.AutoMigrate(&job.Job{}, &job.JobStats{}).GetErrors(); len(errors) != 0 {
		fmt.Printf("Cound not auto migrate tables for reasons below\n %v", errors)
		fmt.Println()
		panic("Could not make/migrate tables")
	}
	Dispatch.SetPersistStorage(db)
	svc := service.New(db, Dispatch)
	ep := endpoints.New(svc)
	r := apischedulerhttp.NewHTTPHandler(ep)
	if verbose || showHTTPDir {
		r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			t, err := route.GetPathTemplate()
			if err != nil {
				return err
			}
			// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
			// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
			p, err := route.GetPathRegexp()
			if err != nil {
				return err
			}
			m, err := route.GetMethods()
			if err != nil {
				return err
			}
			fmt.Println(strings.Join(m, ","), t, p)
			return nil
		})
	}
	log.Infof("Starting Server on port %d", port)
	Dispatch.AddPendingJobs()
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 7 * time.Second,
		Addr:         parsedPort,
		Handler:      r,
	}
	return server.ListenAndServe()
}
