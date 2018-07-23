//Copyright © 2018 Fuf
package config

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DBConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	driver   string
}

var (
	cfg DBConfig
)

func InitVars(cfgType, driver string) error {

	setupLogrusConfig()
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
		return err
	}
	log.WithFields(log.Fields{
		"config": viper.ConfigFileUsed(),
	}).Info("Use")

	setupDBConfig(cfgType, driver)

	return nil
}

func DBConnString() string {
	switch cfg.driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.host, cfg.port, cfg.user, cfg.password, cfg.dbname)
	case "mssql":
		return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d",
			cfg.host, cfg.user, cfg.password, cfg.port)
	case "mysql":
		//user:password@tcp(host:3306)/dbname?charset=utf8
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&interpolateParams=true",
			cfg.user, cfg.password, cfg.host, cfg.port, cfg.dbname)
	default:

	}
	return ""

}

func HttpPort() string {
	return viper.GetString("http_port")
}

func DBDriver() string {
	return cfg.driver
}

func setupLogrusConfig() {
	logName := flag.String("logname", "server.log", "log file name")
	logLevel := flag.String("loglvl", "info", "log level can be: info, debug, error")
	logType := flag.String("logtype", "pro", "log has three type: pro (to JSON), dev (to TTY), debug (to text file)")
	flag.Parse()
	fmt.Printf("flags: %s, %s, %s", *logName, *logType, *logLevel)

	setupLogFormatter(*logType)
	lvl, err := log.ParseLevel(*logLevel)
	if err == nil {
		log.SetLevel(lvl)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if *logType == "dev" {
		return
	}

	file, err := os.OpenFile(*logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to open", *logName, " :", err)
	}

	log.SetOutput(file)

	/*setupLogFormatter(*logType)

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
		FullTimestamp:   true,
		//	DisableTimestamp: false,
	})
	if os.Getenv("LOG_LEVEL") != "" {
		l := os.Getenv("LOG_LEVEL")
		lvl, err := log.ParseLevel(l)
		if err == nil {
			log.SetLevel(lvl)
		}
	} else {
		log.SetLevel(log.InfoLevel)
	}
	*/
}

func setupLogFormatter(logType string) {
	switch logType {
	case "pro":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
		})
	case "dev":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			FullTimestamp:   true,
			//	DisableTimestamp: false,
		})
	case "debug":
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			FullTimestamp:   true,
			//	DisableTimestamp: false,
		})

	}
}

func setupDBConfig(cfgType, driver string) {
	log.Debugf("setupDBConfig type=%s, driver=%s", cfgType, driver)
	switch cfgType {
	case "pro":
		proCfg(driver)
	case "test":
		testCfg(driver)
	default:
		log.Warning("Unknown configuration: ", cfgType)
	}
}

func dbDriver() string {
	ret := os.Getenv("DB_DRIVER")
	if ret == "" {
		ret = viper.GetString("db_driver")
	}
	log.Info("Using database driver = ", ret)
	return ret
}

func proCfg(driver string) {
	if driver == "" {
		//db driver can be declared inside Dockerfile or inside config.yaml
		driver = dbDriver()
	}
	cfg = DBConfig{
		host:     viper.GetString(driver + ".pro_db.host"),
		port:     viper.GetInt(driver + ".pro_db.port"),
		user:     viper.GetString(driver + ".pro_db.user"),
		password: viper.GetString(driver + ".pro_db.pass"),
		dbname:   viper.GetString(driver + ".pro_db.dbname"),
		driver:   driver,
	}
	if os.Getenv("PG_HOST") != "" {
		cfg.host = os.Getenv("PG_HOST")
	}
	if os.Getenv("MYSQL_HOST") != "" {
		cfg.host = os.Getenv("MYSQL_HOST")
	}
	if os.Getenv("MSSQL_HOST") != "" {
		cfg.host = os.Getenv("MSSQL_HOST")
	}

	log.Debug("pro config:", cfg)
}

func testCfg(driver string) {
	if driver == "" {
		driver = viper.GetString("db_driver")
	}
	cfg = DBConfig{
		host:     viper.GetString(driver + ".test_db.host"),
		port:     viper.GetInt(driver + ".test_db.port"),
		user:     viper.GetString(driver + ".test_db.user"),
		password: viper.GetString(driver + ".test_db.pass"),
		dbname:   viper.GetString(driver + ".test_db.dbname"),
		driver:   driver,
	}
	log.Debug("test config:", cfg)
}
