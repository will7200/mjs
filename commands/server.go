package commands

import (
	"bufio"
	"fmt"
	stdlog "log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	colorable "github.com/mattn/go-colorable"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/will7200/mjs/api"
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
	detach            bool
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
	servercmd.Flags().BoolVar(&detach, "detach", false, "detach from console")
	viper.BindPFlag("port", servercmd.Flags().Lookup("port"))
	viper.BindPFlag("database.dbname", servercmd.Flags().Lookup("dbname"))
	viper.BindPFlag("database.connection", servercmd.Flags().Lookup("connection"))
	viper.BindPFlag("interface.workers", servercmd.Flags().Lookup("workers"))
}

func server(cmd *cobra.Command, args []string) error {
	if detach {
		return rundetach(args)
	}
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	var parsedPort string
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
	if err != nil {
		panic(fmt.Sprintf("failed to connect database \ntype %s with connection %s", viper.GetString("dbname"), viper.GetString("connection")))
	}
	if verbose {
		db.LogMode(true)
	}
	Dispatch.SetPersistStorage(db)
	fmt.Println(viper.GetString("database.dbname"))
	if result := db.AutoMigrate(&job.Job{}, &job.JobStats{}).GetErrors(); len(result) != 0 {
		log.Fatal("Couldn't migrate the needed tables shutting down with the following errors:\n", result)
	}
	if filelocation := viper.GetString("logging.filelocation"); filelocation != "" {
		path := filelocation
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			//rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
		)
		if err != nil {
			log.Fatal("failed to create rotatelogs:", err)
		}
		log.AddHook(lfshook.NewHook(lfshook.WriterMap{
			log.InfoLevel:  writer,
			log.ErrorLevel: writer,
		}))

	}
	log.Infof("Starting Server on port %d", port)
	Dispatch.AddPendingJobs()
	return api.StartServer(parsedPort, Dispatch, db)
}

func rundetach(args []string) error {
	for i, v := range os.Args {
		if v == "--detach" {
			os.Args = append(os.Args[:i], os.Args[i+1:len(os.Args)]...)
		}
	}
	cmd := exec.Command(os.Args[0], os.Args[1:len(os.Args)]...)
	file, err := os.Create("log_mjs.txt")
	if err != nil {
		return err
	}
	cmd.Stdout = file
	cmd.Stderr = file
	err = cmd.Start()
	if err != nil {
		fmt.Println("Process exiting with error ", err)
	}
	done := make(chan error)
	timeout := make(chan bool)
	time.AfterFunc(3*time.Second, func() {
		timeout <- true
	})
	//scanner := make(chan bool)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		log.Println("Process exiting with error ", err)
		file, _ = os.Open("log_mjs.txt")
		fs := bufio.NewScanner(file)
		for fs.Scan() {
			text := fs.Text()
			if strings.HasPrefix(text, "Error") {
				log.Println(text)
				break
			}
		}
		return nil
	case _ = <-timeout:
		return nil
	}
}
