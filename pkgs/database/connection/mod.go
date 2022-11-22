package connection

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewConnection(log *zap.Logger, sch *string) *gorm.DB {
	config := getPgConfig()
	connection := config.ToConnectionURI()
	s := "public"
	if sch == nil {
		sch = &s
	}

	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   *sch + ".", // schema name
			SingularTable: false,
		},
		Logger: logger.Default,
	})
	if err != nil {
		log.Error(err.Error())
	}
	return db
}

type PgConfig struct {
	URL         string `env:"DATABASE_URL"`
	Host        string `env:"DATABASE_HOST"`
	Port        int32  `env:"DATABASE_PORT"`
	User        string `env:"DATABASE_USER"`
	Password    string `env:"DATABASE_PASSWORD"`
	Database    string `env:"DATABASE_DB"`
	Environment string `env:"GIN_MODE"`
	AppName     string `env:"APP_NAME"`
}

func (c PgConfig) ToConnectionURI() string {
	return fmt.Sprintf(`postgresql://%s:%s@%s:%v/%s?application_name=%s`, c.User, c.Password, c.Host, c.Port, c.Database, c.AppName)
}

func getPgConfig() *PgConfig {
	c := PgConfig{}
	if err := env.Parse(&c); err != nil {
		log.Fatal(err)
	}
	return &c
}
