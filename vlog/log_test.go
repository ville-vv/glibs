package vlog

import (
	"github.com/golang/mock/gomock"
	"testing"
	"tstl/mock"
)

func TestLogD(t *testing.T) {
	wan := "aa"
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	mockLog := mock.NewMockILogger(ctl)
	gomock.InOrder(
		mockLog.EXPECT().LogD(wan, wan),
		mockLog.EXPECT().LogE(wan, wan),
		mockLog.EXPECT().LogI(wan, wan),
		mockLog.EXPECT().LogW(wan, wan))
	SetLogger(mockLog)
	LogD(wan, wan)
	LogE(wan, wan)
	LogI(wan, wan)
	LogW(wan, wan)
}

func TestDefaultLogger(t *testing.T) {
	DefaultLogger()
}
