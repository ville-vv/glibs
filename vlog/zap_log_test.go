package vlog

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name string
		want *ZapLogger
		cnf  *LogCnf
	}{
		{name: "ss", cnf: &LogCnf{}},
	}
	for _, tt := range tests {
		got := NewZapLogger(tt.cnf)
		tt.want = got
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. NewZapLogger() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestZapLogger_LogD(t *testing.T) {
	type fields struct {
		infoLog *zap.Logger
		logCnf  *LogCnf
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "", args: args{format: "%s++++%s=====%s", args: []interface{}{"sdfsdfs", "eon", "4645"}}, fields: fields{logCnf: &LogCnf{OutPutFile: []string{"stdout"}}}},
	}
	for _, tt := range tests {
		l := NewZapLogger(tt.fields.logCnf)
		l.LogD(tt.args.format, tt.args.args...)
	}
}

func TestZapLogger_LogE(t *testing.T) {
	type fields struct {
		infoLog *zap.Logger
		accLog  *zap.Logger
		logCnf  *LogCnf
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		l := NewZapLogger(tt.fields.logCnf)
		l.LogE(tt.args.format, tt.args.args...)
	}
}

func TestZapLogger_LogW(t *testing.T) {
	type fields struct {
		infoLog *zap.Logger
		accLog  *zap.Logger
		logCnf  *LogCnf
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		l := NewZapLogger(tt.fields.logCnf)
		l.LogW(tt.args.format, tt.args.args...)
	}
}

func TestZapLogger_LogI(t *testing.T) {
	type fields struct {
		infoLog *zap.Logger
		accLog  *zap.Logger
		logCnf  *LogCnf
	}
	type args struct {
		format string
		args   []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		l := NewZapLogger(tt.fields.logCnf)
		l.LogI(tt.args.format, tt.args.args...)
	}
}
