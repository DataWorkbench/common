package gormwrap

import (
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

const (
	callBackBeforeName = "OpenTracing:Before"
	callBackAfterName  = "OpenTracing:After"
)

var traceComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "gorm"}

type openTracingPlugin struct {
	tracer opentracing.Tracer
}

func newOpenTracingPlugin(tracer opentracing.Tracer) gorm.Plugin {
	return &openTracingPlugin{tracer: tracer}
}

// Implements gorm.Plugin
func (pl *openTracingPlugin) Name() string {
	return "OpenTracingPlugin"
}

// Implements gorm.Plugin
func (pl *openTracingPlugin) Initialize(db *gorm.DB) (err error) {
	// Register before function.
	err = db.Callback().Create().Before("gorm:CreateBefore").Register(callBackBeforeName, pl.before)
	if err != nil {
		return
	}
	err = db.Callback().Query().Before("gorm:QueryBefore").Register(callBackBeforeName, pl.before)
	if err != nil {
		return
	}
	err = db.Callback().Delete().Before("gorm:DeleteBefore").Register(callBackBeforeName, pl.before)
	if err != nil {
		return
	}
	err = db.Callback().Update().Before("gorm:UpdateBefore").Register(callBackBeforeName, pl.before)
	if err != nil {
		return
	}
	err = db.Callback().Row().Before("gorm:RowBefore").Register(callBackBeforeName, pl.before)
	if err != nil {
		return
	}
	err = db.Callback().Raw().Before("gorm:RawBefore").Register(callBackBeforeName, pl.before)
	if err != nil {
		return
	}

	// Register after function.
	err = db.Callback().Create().After("gorm:CreateAfter").Register(callBackAfterName, pl.after)
	if err != nil {
		return
	}
	err = db.Callback().Query().After("gorm:QueryAfter").Register(callBackAfterName, pl.after)
	if err != nil {
		return
	}
	err = db.Callback().Delete().After("gorm:DeleteAfter").Register(callBackAfterName, pl.after)
	if err != nil {
		return
	}
	err = db.Callback().Update().After("gorm:UpdateAfter").Register(callBackAfterName, pl.after)
	if err != nil {
		return
	}
	err = db.Callback().Row().After("gorm:RowAfter").Register(callBackAfterName, pl.after)
	if err != nil {
		return
	}
	err = db.Callback().Raw().After("gorm:RawAfter").Register(callBackAfterName, pl.after)
	if err != nil {
		return
	}
	return
}

func (pl *openTracingPlugin) before(db *gorm.DB) {
	var parentCtx opentracing.SpanContext

	ctx := db.Statement.Context

	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	sql := db.Statement.SQL.String()
	opName := strings.Split(sql, " ")[0]

	span := pl.tracer.StartSpan(
		db.Name()+opName,
		opentracing.ChildOf(parentCtx),
		ext.SpanKindRPCClient,
		traceComponentTag,
		opentracing.Tag{Key: string(ext.DBType), Value: db.Name()},
	)

	span.LogFields(tracerLog.String("sql", db.Dialector.Explain(sql, db.Statement.Vars...)))

	db.Statement.Context = opentracing.ContextWithSpan(ctx, span)
}

func (pl *openTracingPlugin) after(db *gorm.DB) {
	span := opentracing.SpanFromContext(db.Statement.Context)

	// Error
	if err := db.Error; err != nil {
		ext.Error.Set(span, true)
		span.LogFields(tracerLog.Error(err))
	}
	span.Finish()
}
