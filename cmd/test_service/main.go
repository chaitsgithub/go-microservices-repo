package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"chaits.org/microservices-repo/pkg/errors"
	"chaits.org/microservices-repo/pkg/general/config"
	"chaits.org/microservices-repo/pkg/general/logger"
	"chaits.org/microservices-repo/pkg/network/httpclient"
	sqldb "chaits.org/microservices-repo/pkg/storage/sqldb/connectors"
	constants "chaits.org/microservices-repo/test_service/internal"
)

var appLogger *logger.Logger
var appErr *errors.AppError

func main() {
	appLogger = logger.NewLogger(constants.SERVICE_NAME, constants.ENVIRONMENT)
	var utilToTest string

	if os.Args[1] == "" {
		utilToTest = "All"
	} else {
		utilToTest = os.Args[1]
	}

	switch true {
	case utilToTest == "logger":
		testLoggerUtility()
	case utilToTest == "config":
		fmt.Println()
		testConfigUtility()
	case utilToTest == "errors":
		fmt.Println()
		testErrorsUtility()
	case utilToTest == "httpclient":
		fmt.Println()
		testHttpClient()
	case utilToTest == "db":
		testDB()
	}
}

type Student struct {
	id        int
	firstName string
	lastName  string
	age       int
}

// Test DB
func testDB() {
	log.Println("Testing DB Utility")

	dbConnector, err := sqldb.NewConnector(sqldb.DB_POSTGRES)
	if err != nil {
		log.Fatalf("Error getting DB Connector. Error : %w", err)
	}

	dbConn, err := dbConnector.Connect()
	if err != nil {
		log.Fatalf("Error connecting to DB. Error : %w", err)
	}
	defer dbConn.Close()

	rows, err := dbConn.Query("select * from students")
	if err != nil {
		log.Fatalf("Error querying table. Error : %w", err)
	}

	for rows.Next() {
		var s Student
		rows.Scan(&s.id, &s.firstName, &s.lastName, &s.age)
		log.Println(s)
	}
}

// Test HTTP Client Utility
func testHttpClient() {

	log.Println("*** Testing HttpClient Utility")
	httpClient := httpclient.NewHTTPClient(httpclient.WithTimeout(5*time.Second), httpclient.WithRetry(3, 1*time.Second, http.StatusNotFound))
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://httpbin.org/get", nil)
	if err != nil {
		log.Fatalln("Unable to create http request")
	}

	resp, err := httpClient.Do(context.Background(), req)
	var logMessage string
	if err != nil {
		log.Printf("Error : %v", err)
		logMessage = "Error Processing Request"
	} else {
		logMessage = "Request Processed Successfully!"
	}

	appLogger.LogHttpMessage(logMessage, resp, logger.ContextData{})
}

// Code to Test Logger Utility
func testLoggerUtility() {

	log.Println("*** Testing Logger Utility")
	appLogger.LogMessage(fmt.Sprintf("%s ran successfully #1", constants.SERVICE_NAME), logger.LogParms{})
	// appLogger.LogMessage(fmt.Sprintf("%s ran successfully #2", constants.SERVICE_NAME), logger.LogParms{})
	// appLogger.LogMessage(fmt.Sprintf("%s ran successfully #3", constants.SERVICE_NAME), logger.LogParms{})

}

func testConfigUtility() {
	log.Println("*** Testing Config Utility")
	configs := config.InitConfigs("dev")
	configs.LoadConfigs()

	log.Println("devconfig#1 : ", configs.GetConfig("devconfig#1"))
	// log.Println(configs.GetConfig("devconfig#3"))
	// log.Println("All Configs: ")
	// configs.PrintAllKeys()

}

func testErrorsUtility() {

	log.Println("*** Testing Errors Utility")
	err := testErrors()
	if err != nil {
		if errors.As(err, &appErr) {
			log.Println("Error Code : ", appErr.Code)
			log.Println("Error Message : ", appErr.Message)
		}
		log.Println(err)
	}

}

func testErrors() error {
	err := doSomething()
	if err != nil {
		return errors.Wrap("HANDLER_ERR", "Processing failed", err)
	}
	return nil
}

func doSomething() error {
	return errors.New("DB_ERR", "Failed to connect to database")
}
