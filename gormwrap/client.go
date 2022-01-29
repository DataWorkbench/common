package gormwrap

import (
	"context"
	"database/sql"
	"fmt"
	"go/types"
	"strings"
	"time"

	"github.com/DataWorkbench/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/DataWorkbench/common/gtrace"
)

const (
	CondEqual         = " = ? "
	CondNotEqual      = " <> ? "
	CondIn            = " IN ? "
	CondLike          = " LIKE ? "
	CondLikePlacehold = "%"

	JointAnd   = " AND "
	ReverseFmt = "%s DESC"
)

type MySQLConfig struct {
	// Hosts sample "127.0.0.1:3306,127.0.0.1:3307,127.0.0.1:3308"
	Hosts       string `json:"hosts"         yaml:"hosts"         env:"HOSTS"                     validate:"required"`
	Users       string `json:"users"         yaml:"users"         env:"USERS"                     validate:"required"`
	Password    string `json:"password"      yaml:"password"      env:"PASSWORD"                  validate:"required"`
	Database    string `json:"database"      yaml:"database"      env:"DATABASE"                  validate:"required"`
	MaxIdleConn int    `json:"max_idle_conn" yaml:"max_idle_conn" env:"MAX_IDLE_CONN,default=16"  validate:"required"`
	MaxOpenConn int    `json:"max_open_conn" yaml:"max_open_conn" env:"MAX_OPEN_CONN,default=128" validate:"required"`
	// gorm log level: 1 => Silent, 2 => Error, 3 => Warn, 4 => Info
	LogLevel int `json:"log_level" yaml:"log_level" env:"LOG_LEVEL,default=4" validate:"gte=1,lte=4"`
	// ConnMaxLifetime unit seconds
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime" env:"CONN_MAX_LIFETIME,default=10m" validate:"required"`
	// SlowThreshold time 0 indicates disabled
	SlowThreshold time.Duration `json:"slow_threshold" yaml:"slow_threshold" env:"SLOW_THRESHOLD,default=2s" validate:"gte=0"`
}

// NewMySQLConn return a grom.DB by mysql driver
// NOTICE: Must set glog.Logger into the ctx by glow.WithContext
func NewMySQLConn(ctx context.Context, cfg *MySQLConfig) (db *gorm.DB, err error) {
	lp := glog.FromContext(ctx)

	defer func() {
		if err != nil {
			lp.Error().Error("gorm: create mysql connection error", err).Fire()
			db = nil
		}
	}()

	lp.Info().Msg("gorm: connecting to mysql").String("hosts", cfg.Hosts).String("database", cfg.Database).Fire()

	hosts := strings.Split(strings.ReplaceAll(cfg.Hosts, " ", ""), ",")
	if len(hosts) == 0 {
		err = fmt.Errorf("invalid hosts %s", cfg.Hosts)
		return
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Users, cfg.Password, hosts[0], cfg.Database,
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//SkipDefaultTransaction: true,
		Logger: &Logger{
			Level:         LogLevel(cfg.LogLevel),
			SlowThreshold: cfg.SlowThreshold,
			Output:        lp,
		},
	})
	if err != nil {
		return
	}

	// Set connection pool
	var sqlDB *sql.DB
	if sqlDB, err = db.DB(); err != nil {
		return
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	tracer := gtrace.TracerFromContext(ctx)
	if err = db.Use(newOpenTracingPlugin(tracer)); err != nil {
		return
	}

	// TODO: Adds multiple databases if necessary
	// import gorm.io/plugin/dbresolver
	return
}

type Condition struct {
	// where cond: k = v
	Values    map[string]interface{}
	Operators map[string]string

	Offset  int
	Limit   int
	Order   string
	Reverse bool
}

// default operator is CondEqual
func (c *Condition) Update(column string, value interface{}, operator ...string) {
	c.Values[column] = value
	if len(operator) > 0 {
		c.Operators[column] = operator[0]
	} else {
		switch value.(type) {
		case types.Slice:
			c.Operators[column] = CondIn
		case types.Array:
			c.Operators[column] = CondIn
		default:
			c.Operators[column] = CondEqual
		}
	}
}

func (c *Condition) UpdateLimit(offset, limit int) {
	c.Offset = offset
	c.Limit = limit
}

func (c *Condition) UpdateOrder(order string, reverse bool) {
	c.Order = order
	c.Reverse = reverse
}

// conditions: column: values(string, int or slice)
func (c *Condition) Build(tx *gorm.DB) *gorm.DB {
	var q, joint string
	var a []interface{}
	for column, v := range c.Values {
		q += joint + column + c.Operators[column]
		// handle like
		if c.Operators[column] == CondLike {
			a = append(a, fmt.Sprintf("%s%v%s", CondLikePlacehold, v, CondLikePlacehold))
		} else {
			a = append(a, v)
		}
		joint = JointAnd
	}

	if q != "" {
		tx = tx.Where(q, a...)
	}
	if c.Limit > 0 {
		tx = tx.Offset(c.Offset).Limit(c.Limit)
	}
	if c.Order != "" {
		if c.Reverse {
			tx = tx.Order(fmt.Sprintf(ReverseFmt, c.Order))
		} else {
			tx = tx.Order(c.Order)
		}
	}
	return tx
}
