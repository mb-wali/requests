package main

import (
	"github.com/cyverse-de/echo-middleware/redoc"
	"github.com/cyverse-de/requests/api"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

// buildLoggerEntry sets some logging options then returns a logger entry with some custom fields
// for convenience.
func buildLoggerEntry() *logrus.Entry {

	// Enable logging the file name and line number.
	logrus.SetReportCaller(true)

	// Set the logging format to JSON for now because that's what Echo's middleware uses.
	logrus.SetFormatter(&logrus.JSONFormatter{})

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
	serviceInfo, err := getSwaggerServiceInfo()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Define the API.
	a := api.API{
		Echo:    e,
		Title:   serviceInfo.Title,
		Version: serviceInfo.Version,
	}

	// Define the API endpoints.
	e.GET("/", a.RootHandler)

	// Start the service.
	e.Logger.Fatal(e.Start(":8080"))
}
