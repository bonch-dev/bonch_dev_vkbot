package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	sentry "github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql"
	env "github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"vk-bot/handlers"

	"github.com/SevereCloud/vksdk/5.92/api"
	"github.com/SevereCloud/vksdk/5.92/callback"
)

var vkapi api.VK
var cb callback.Callback
var logger *log.Logger

func init() {
	if err := env.Load(); err != nil {
		if _, err := os.Stat("./.env.example"); os.IsExist(err) {
			log.Fatal("File .env.example exists! Fill it and rename to .env, before we go")
		}
		log.Fatalf("Error .env loading: %v", err.Error())
	} else {
		log.Info(".env file loaded")
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://f6d24d296ccc451a9c3b3ca7441d8774@sentry.io/1542414",
	}); err != nil {
		log.Error("Error in Sentry init")
	}

	logger = log.New()

	if os.Getenv("DEBUG") == "true" {
		logger.SetLevel(log.TraceLevel)
	} else {
		logger.SetLevel(log.ErrorLevel)
	}

	vkapi = api.Init(os.Getenv("TOKEN"))

	limit, err := strconv.ParseInt(os.Getenv("API_LIMIT"), 10, 32)
	if err != nil {
		limit = 0
	}

	vkapi.Limit = int(limit)

	cb.ConfirmationKey = os.Getenv("CONFIRMATION_KEY")
}

func main() {
	db := getDBClient(logger)

	defHandler := handlers.NewDefHandler(vkapi, db, logger)

	cb.MessageNew(handlers.TextHandle(defHandler))

	http.HandleFunc("/", cb.HandleFunc)

	log.Info("Bot started at :4043")
	sentry.CaptureMessage("Bot started at :4043")
	if err := http.ListenAndServe(":4043", nil); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}

	sentry.Flush(time.Second * 5)
}

func getDBClient(logger *log.Logger) *sql.DB {
	var (
		clientConfig string
		driver       = os.Getenv("DRIVER")
		driverHost   = os.Getenv("DRIVER_HOST")
		driverUser   = os.Getenv("DRIVER_USER")
		driverPass   = os.Getenv("DRIVER_PASS")
		driverDb     = os.Getenv("DRIVER_DB")
	)

	switch driver {
	case "mysql":
		clientConfig = fmt.Sprintf("%v:%v@tcp(%v)/%v", driverUser, driverPass, driverHost, driverDb)
	default:
		clientConfig = fmt.Sprintf("%v:%v@%v/%v", driverUser, driverPass, driverHost, driverDb)
	}

	logger.Infof("Connecting to DB with conf: %v", clientConfig)
	sentry.CaptureMessage(fmt.Sprintf("Connecting to DB with conf: %v", clientConfig))

	client, e := sql.Open(driver, clientConfig)
	if e != nil {
		sentry.CaptureException(errors.WithMessage(e, fmt.Sprintf("Can't load DB with this data: %v", clientConfig)))
		sentry.Flush(time.Second * 5)
		logger.Errorln(e)
		logger.Fatalf("Can't load DB with this data: %v", clientConfig)
	}

	if err := client.Ping(); err != nil {
		logger.Fatal(err)
	}

	return client
}
