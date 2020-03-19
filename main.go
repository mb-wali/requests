package main

import (
	"flag"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/cyverse-de/configurate"
	"github.com/cyverse-de/echo-middleware/redoc"
	"github.com/cyverse-de/requests/api"
	"github.com/cyverse-de/requests/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

var (
	cfgPath = flag.String("config", "/etc/iplant/de/jobservices.yml", "The path to the config file")
	port    = flag.String("port", "8080", "The port to listen to")
	debug   = flag.Bool("debug", false, "Enable debug logging")
)

func init() {
	flag.Parse()
}

// buildLoggerEntry sets some logging options then returns a logger entry with some custom fields
// for convenience.
func buildLoggerEntry() *logrus.Entry {

	// Enable logging the file name and line number.
	logrus.SetReportCaller(true)

	// Set the logging format to JSON for now because that's what Echo's middleware uses.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Enable debugging if we're supposed to.
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Return the custom log entry.
	return logrus.WithFields(logrus.Fields{
		"service": "requests",
		"art-id":  "requests",
		"group":   "org.cyverse",
	})
}

func main() {
	e := echo.New()

	// Set a custom logger.
	e.Logger = Logger{Entry: buildLoggerEntry()}

	// Add middleware.
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(redoc.Serve(redoc.Opts{Title: "DE Administrative Requests API Documentation"}))

	// Load the service information from the Swagger JSON.
	e.Logger.Info("loading service information")
	serviceInfo, err := getSwaggerServiceInfo()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Load the configuration file.
	e.Logger.Info("loading the configuration file")
	cfg, err := configurate.Init(*cfgPath)
	if err != nil {
		e.Logger.Fatalf("unable to load the configuration file: %s", err.Error())
	}

	// Initialize the database connection.
	e.Logger.Info("establishing the database connection")
	databaseURI := cfg.GetString("db.uri")
	db, err := db.InitDatabase("postgres", databaseURI)
	if err != nil {
		e.Logger.Fatalf("service initialization failed: %s", err.Error())
	}

	// Define the API.
	a := api.API{
		Echo:    e,
		Title:   serviceInfo.Title,
		Version: serviceInfo.Version,
		DB:      db,
	}

	// Define the API endpoints.
	e.GET("/", a.RootHandler)
	e.GET("/request-types", a.GetRequestTypesHandler)
	e.POST("/request-types/:name", a.RegisterRequestTypeHandler)
	e.GET("/request-types/:name", a.GetRequestTypeHandler)
	e.GET("/request-status-codes", a.GetRequestStatusCodesHandler)

	// Start the service.
	e.Logger.Info("starting the service")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
