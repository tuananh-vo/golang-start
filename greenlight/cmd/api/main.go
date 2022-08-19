package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"greenlight/internal/data"
	"greenlight/internal/jsonlog"
	"greenlight/internal/mailer"

	_ "github.com/lib/pq"
)

// Remove the hardcoded version number and make version a variable instead of a constant.
var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	// Add a new limiter struct containing fields for the requests-per-second and burst
	// values, and a boolean field which we can use to enable/disable rate limiting
	// altogether.
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}

	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}

	// Add a cors struct and trustedOrigins field with the type []string.
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	// Declare an instance of the config struct.
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Use the value of the GREENLIGHT_DB_DSN environment variable as the default value
	// for our db-dsn command-line flag.
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	// flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:pass@localhost/greenlight", "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice the default values that we're using?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	// Create command line flags to read the setting values into the config struct.
	// Notice that we use true as the default for the 'enabled' setting?
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// Read the SMTP server configuration settings into the config struct, using the
	// Mailtrap settings as the default values. IMPORTANT: If you're following along,
	// make sure to replace the default values for smtp-username and smtp-password
	// with your own Mailtrap credentials.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "89d576bae9d48e", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "fb6f87e0af1352", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@anhvt>", "SMTP sender")

	// Use the flag.Func() function to process the -cors-trusted-origins command line
	// flag. In this we use the strings.Fields() function to split the flag value into a
	// slice based on whitespace characters and assign it to our config struct.
	// Importantly, if the -cors-trusted-origins flag is not present, contains the empty
	// string, or contains only whitespace, then strings.Fields() will return an empty
	// []string slice.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	// Create a new version boolean flag with the default value of false.
	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	// If the version flag value is true, then print out the version number and
	// immediately exit.
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		// Print out the contents of the buildTime variable.
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	// Initialize a new jsonlog.Logger which writes any messages *at or above* the INFO
	// severity level to the standard out stream.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.PrintInfo("database connection pool established", nil)

	// Publish a new "version" variable in the expvar handler containing our application
	// version number (currently the constant "1.0.0").
	expvar.NewString("version").Set(version)

	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))
	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	// Declare an instance of the application struct, containing the config struct and
	// the logger.
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	// Declare a HTTP server with some sensible timeout settings, which listens on the
	// port provided in the config struct and uses the servemux we created above as the
	// handler.
	// srv := &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", cfg.port),
	// 	Handler:      app.routes(),
	// 	ErrorLog:     log.New(logger, "", 0),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 30 * time.Second,
	// }
	// // Start the HTTP server.
	// logger.PrintInfo("starting server", map[string]string{
	// 	"addr": srv.Addr,
	// 	"env":  cfg.env,
	// })
	// err = srv.ListenAndServe()
	// logger.PrintFatal(err, nil)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	// Return the sql.DB connection pool.
	return db, nil
}
