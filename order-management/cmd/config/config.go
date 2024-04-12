package config

import "time"

// Global defines the global configuration values
var Global = struct {
	RabbitmqScheme   string
	RabbitmqUsername string
	RabbitmqPassword string
	RabbitmqAddress  string
	RabbitmqPort     uint16
	RabbitmqVhost    string

	PostgresAddress  string
	PostgresUsername string
	PostgresPassword string
	PostgresDatabase string

	Delay time.Duration
}{

	RabbitmqScheme:   "amqp",
	RabbitmqUsername: "",
	RabbitmqPassword: "",
	RabbitmqAddress:  "",
	RabbitmqPort:     5672,
	RabbitmqVhost:    "/",

	PostgresAddress:  "postgres",
	PostgresUsername: "postgres",
	PostgresPassword: "postgres",
	PostgresDatabase: "products",
}
