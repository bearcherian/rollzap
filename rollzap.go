package rollzap

import (
	"encoding/json"
	"log"

	rollbar "github.com/rollbar/rollbar-go"
	"go.uber.org/zap/zapcore"
)

type levelEnabler struct {
	minLevel zapcore.Level
}

// RollbarCore is a custom core to send logs to Rollbar. Add the core using zapcore.NewTee
type RollbarCore struct {
	levelEnabler
	coreFields map[string]interface{}
}

// NewRollbarCore creates a new core to transmit logs to rollbar. rollbar token and other options should be set before creating a new core
func NewRollbarCore(minLevel zapcore.Level) *RollbarCore {

	return &RollbarCore{
		levelEnabler: levelEnabler{
			minLevel: minLevel,
		},
		coreFields: make(map[string]interface{}),
	}
}

func (le *levelEnabler) Enabled(l zapcore.Level) bool {
	return l >= le.minLevel
}

// With provides structure
func (c *RollbarCore) With(fields []zapcore.Field) zapcore.Core {

	fieldMap := fieldsToMap(fields)

	for k, v := range fieldMap {
		c.coreFields[k] = v
	}

	return c
}

// Check determines if this should be sent to roll bar based on LevelEnabler
func (c *RollbarCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.levelEnabler.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, c)
	}
	return checkedEntry
}

func (c *RollbarCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {

	fieldMap := fieldsToMap(fields)

	if len(c.coreFields) > 0 {
		if coreFieldsMap, err := json.Marshal(c.coreFields); err != nil {
			log.Println("Unable to parse json for coreFields")
		} else {
			fieldMap["coreFields"] = string(coreFieldsMap)
		}
	}

	if entry.LoggerName != "" {
		fieldMap["logger"] = entry.LoggerName
	}
	if entry.Caller.TrimmedPath() != "" {
		fieldMap["file"] = entry.Caller.TrimmedPath()
	}

	switch entry.Level {
	case zapcore.DebugLevel:
		rollbar.Debug(entry.Message, fieldMap)
	case zapcore.InfoLevel:
		rollbar.Info(entry.Message, fieldMap)
	case zapcore.WarnLevel:
		rollbar.Warning(entry.Message, fieldMap)
	case zapcore.ErrorLevel:
		rollbar.Error(entry.Message, fieldMap)
	case zapcore.DPanicLevel:
		rollbar.Critical(entry.Message, fieldMap)
	case zapcore.PanicLevel:
		rollbar.Critical(entry.Message, fieldMap)
	case zapcore.FatalLevel:
		rollbar.Critical(entry.Message, fieldMap)

	}
	return nil
}

// Sync flushes
func (c *RollbarCore) Sync() error {
	rollbar.Wait()
	return nil
}

func fieldsToMap(fields []zapcore.Field) map[string]interface{} {
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	m := make(map[string]interface{})
	for k, v := range enc.Fields {
		m[k] = v
	}
	return m
}
