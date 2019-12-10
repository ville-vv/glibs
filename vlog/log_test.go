package vlog

import (
	"testing"
)

func TestLogD(t *testing.T) {
	//wan := "aa"
	//ctl := gomock.NewController(t)
	//defer ctl.Finish()
	//mockLog := mock.NewMockILogger(ctl)
	//gomock.InOrder(
	//	mockLog.EXPECT().LogD(wan, wan),
	//	mockLog.EXPECT().LogE(wan, wan),
	//	mockLog.EXPECT().LogI(wan, wan),
	//	mockLog.EXPECT().LogW(wan, wan))
	//SetLogger(mockLog)
	//LogD(wan, wan)
	//LogE(wan, wan)
	//LogI(wan, wan)
	//LogW(wan, wan)
}

func TestDefaultLogger(t *testing.T) {
	DefaultLogger()
	LogD("13465")
	LogD("1346sdfsddddddddddddddddddddd5")
	LogD("sdfsdsdfsdafsdfsdfsdfsd")
	LogD("sdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsdd")
	LogD("sdfsdsdfsdafsdfsdsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddsdfsdsdfsdafsdfsdfsdfsddfsdfsd")
	LogD("sdfsdsdfsdafsdfsdfsdfsd")
	LogD("sdfsdsdfsdafsdfsdfsdfsd")

}
func TestSetLogger(t *testing.T) {
	cnf := &LogCnf{
		OutPutErrFile: []string{},
		OutPutFile:    []string{"./log/v_log_record.log", "stdout"},
		ProgramName:   ProgramName,
		Level:         LogLevelDebug,
	}
	SetLogger(NewZapLogger(cnf))
	LogD("13465")
}
