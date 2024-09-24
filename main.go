package main

import (
	"dots-api/bootstrap"
	"dots-api/lib/psql"
	"dots-api/lib/utils"
	"dots-api/services/api"
	"dots-api/services/worker/command"
	"fmt"
	"log"
	"os"

	"path/filepath"
	"runtime"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/urfave/cli/v2"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	config     utils.Config
	debug      = false

	// app the base of skeleton
	app *bootstrap.App
)

// EnvConfigPath environtment variable that set the config path
const EnvConfigPath = "REBEL_CLI_CONFIG_PATH"

// setup initialize the used variable and dependencies
func setup() {
	configFile := os.Getenv(EnvConfigPath)
	if configFile == "" {
		configFile = "./config.json"
	}

	log.Println(configFile)

	config = utils.NewViperConfig(basepath, configFile)

	debug = config.GetBool("app.debug")
	validator := bootstrap.SetupValidator(config)
	cLog := bootstrap.SetupLogger(config)

	// connect to redis cache
	rdCache, err := bootstrap.SetupRedis(
		config.GetString("db.redis.addr"),
		config.GetString("db.redis.password"),
		1,
	)
	if err != nil {
		fmt.Println("[redis-cache] " + err.Error())
	}

	// connect to database
	db, err := psql.Connect(config.GetString("db.psql_dsn"))
	if err != nil {
		panic(err)
	}

	// newRelic
	monitoring, err := newrelic.NewApplication(
		newrelic.ConfigAppName(config.GetString("new_relic.relic_name")),
		newrelic.ConfigLicense(config.GetString("new_relic.license_key")),
	)
	if err != nil {
		panic(err)
	}

	app = &bootstrap.App{
		Debug:     debug,
		Config:    config,
		Validator: validator,
		Log:       cLog,
		DB:        db,
		Redis:     rdCache,
		NewRelic:  monitoring,
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setup()

	cmd := &cli.App{
		Name:  "Verein Core",
		Usage: "Verein Core, cli",
		Commands: []*cli.Command{
			{
				Name:   "api",
				Usage:  "API core service, Run on http 1/1",
				Flags:  api.Flags,
				Action: api.Boot{App: app}.Start,
			},
			{
				Name:   "room-reminder",
				Usage:  "Scheduler service, Run on crontab",
				Action: command.Contract{App: app}.UserRoomReminder,
			},
			{
				Name:   "tournament-reminder",
				Usage:  "Scheduler service, Run on crontab",
				Action: command.Contract{App: app}.UserTournamentReminder,
			},
			{
				Name:   "consumer-user-badges",
				Usage:  "Scheduler service, Run on crontab",
				Action: command.Contract{App: app}.ConsumerUserBadge,
			},
			{
				Name:   "consumer-badges",
				Usage:  "Scheduler service, Run on crontab",
				Action: command.Contract{App: app}.ConsumerBadges,
			},
			{
				Name:   "set-inactive-room-and-tournament",
				Usage:  "Scheduler service, Run on crontab",
				Action: command.Contract{App: app}.UpdateStatusRoomAndTournament,
			},
		},
		Action: func(cli *cli.Context) error {
			fmt.Printf("%s version@%s\n", cli.App.Name, "2.1")
			return nil
		},
	}

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
