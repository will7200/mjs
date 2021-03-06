package commands

import (
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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
	pgrpc "github.com/will7200/mjs/apischeduler/grpc"
	"github.com/will7200/mjs/apischeduler/grpc/pb"
	apischedulerhttp "github.com/will7200/mjs/apischeduler/http"
	"github.com/will7200/mjs/apischeduler/service"
	"github.com/will7200/mjs/job"
	"github.com/will7200/mjs/utils/hooks"
	"github.com/rifflock/lfshook"
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
	startgpc          bool
	dispatcherVerbose bool
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
	servercmd.Flags().BoolVar(&startgpc, "startgpc", false, "Start gpc server")
	servercmd.Flags().Int("grpcport", 4005, "Port that grpc will run on")
	servercmd.Flags().BoolVar(&dispatcherVerbose, "dispatcherVerbose", false, "Show ouput of commands from dispatch workers")
	viper.BindPFlag("interface.port", servercmd.Flags().Lookup("port"))
	viper.BindPFlag("database.dbname", servercmd.Flags().Lookup("dbname"))
	viper.BindPFlag("database.connection", servercmd.Flags().Lookup("connection"))
	viper.BindPFlag("interface.workers", servercmd.Flags().Lookup("workers"))
	viper.BindPFlag("interface.home", servercmd.Flags().Lookup("homedir"))
	viper.BindPFlag("verbose", servercmd.Flags().Lookup("verbose"))
	viper.BindPFlag("interface.grpcport", servercmd.Flags().Lookup("grpcport"))
	viper.SetEnvPrefix("mjs") // will be uppercased automatically
}
func server(cmd *cobra.Command, args []string) error {
	verbose = viper.GetBool("verbose")
	if verbose {
		log.Info("changing to debug")
		log.SetLevel(log.DebugLevel)
	}
	var parsedPort string
	port = viper.GetInt("interface.port")
	if port != 0 {
		parsedPort = fmt.Sprintf(":%d", port)
	} else {
		parsedPort = ":4004"
	}
	log.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp:true})
	log.SetOutput(colorable.NewColorableStdout())
	log.AddHook(hooks.NewHook(&hooks.CallerHookOptions{
		Field: "src",
		Flags: stdlog.Lshortfile,
	}))
	if viper.GetString("logging.infolevel") != "" && viper.GetString("logging.errorlevel") != ""{
		hook := lfshook.NewHook(
			lfshook.PathMap{
				log.InfoLevel: viper.GetString("logging.infolevel"),
				log.ErrorLevel: viper.GetString("logging.errorlevel"),
			},
		)
		hook.SetFormatter(&log.JSONFormatter{})
		log.AddHook(hook)
	}
	Dispatch = &job.Dispatcher{}
	if dispatcherVerbose {
		Dispatch.SetVerbose(true)
	}
	Dispatch.StartDispatcher(viper.GetInt("interface.workers"))
	db, err = gorm.Open(viper.GetString("database.dbname"), viper.GetString("database.connection"))
	if verbose {
		db.LogMode(false)
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
	svc := service.New(db, Dispatch, log.StandardLogger())
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
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 7 * time.Second,
		Addr:         parsedPort,
		Handler:      r,
	}
	if startgpc {
		go setupGRPC(svc)
	}
	Dispatch.AddPendingJobs()
	return server.ListenAndServe()
}

func setupGRPC(svc service.APISchedulerService) {
	ss := pgrpc.NewGRPC(svc)
	port := fmt.Sprintf(":%d", viper.GetInt("interface.grpcport"))
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAPISchedulerServer(s, ss)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Infof("Starting GRPC Server on port %d", viper.GetInt("interface.grpcport"))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

/*
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
*/
