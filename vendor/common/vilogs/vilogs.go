package viLogs

import (
	"fmt"
	"log"
	"os"
	"path"

	logging "github.com/op/go-logging"
)

//日志打印的参数
type LogUserParams struct {
	color int
}

// 日志配置
type ViLogsConfig struct {
	ProjectName string
	FilePath    string
	FileName    string
	LogLeavel   string
	//Params      *LogUserParams
}

//日志管理结构体
type ViLogsManage struct {
	LogInfo      *logging.Logger
	LogConf      *ViLogsConfig
	LogModelName string
	LogHoseName  string
}

const (
	CRITICAL = iota
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

var (
	vilogs *ViLogsManage

	//日志打印的级别
	vilevel = map[string]int{
		"Critical": CRITICAL,
		"ERRO":     ERROR,
		"WARN":     WARNING,
		"NOTI":     NOTICE,
		"INFO":     INFO,
		"DEBU":     DEBUG,
	}
)

func getHostName() (hostname string, err error) {
	hostname, err = os.Hostname()
	return
}

func init() {
	//InitLogs(nil, "vil_tools")
}

//模块初始化
func InitLogs(logconf *ViLogsConfig, logmodel string) {
	if logconf == nil {
		//设置默认配置参数
		conf := setDefaultConfig()
		vilogs = NewViLogsManage(conf, logmodel)

	} else {
		vilogs = NewViLogsManage(logconf, logmodel)
	}

}

func setDefaultConfig() *ViLogsConfig {
	var conf *ViLogsConfig
	conf = new(ViLogsConfig)
	conf.FileName = "viltools_logs.log"
	conf.FilePath = "./logs"
	conf.LogLeavel = "Debug"
	conf.ProjectName = "vil_tools"
	return conf
}

/*******************************************************************************
// @describe : logconf 日志参数配置
// @param	 : logmodel 日志模块名称
// @return   : ViLogsManage 日志管理对象
*******************************************************************************/
func NewViLogsManage(logconf *ViLogsConfig, logmodel string) *ViLogsManage {
	lg := new(ViLogsManage)
	if logconf == nil {
		fmt.Println("logconf is nil")
		//设置默认的日志参数
		lg.LogConf = setDefaultConfig()
	} else {
		lg.LogConf = logconf
	}
	lg.LogModelName = logmodel
	lg.LogHoseName, _ = getHostName()
	lg.SetLogger()
	return lg
}

func (lg *ViLogsManage) SetLogger() {
	//创建该目录
	if err := os.MkdirAll(lg.LogConf.FilePath, 0777); err != nil {
		fmt.Println("mkdir error:", err)
		panic(0)
	}
	lg.LogInfo = lg.setLogBackend()
}

func (lg *ViLogsManage) setLogBackend() *logging.Logger {
	//打开文件
	file, err := os.OpenFile(path.Join(lg.LogConf.FilePath, lg.LogConf.FileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("init log failed", err)
		panic(0)
	}

	//日志打印格式
	format := logging.MustStringFormatter(
		`%{color:bold}%{time:2006-01-02 15:04:05.000} %{level:.4s} %{shortfile} %{shortfunc} %{color:reset}%{message}`,
	)
	var lev *logging.Logger
	//设置日志打印位置为文件
	backend := logging.NewLogBackend(file, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.Level(vilevel[lg.LogConf.LogLeavel]), lg.LogModelName)

	//设置模块名称
	lev = logging.MustGetLogger(lg.LogModelName)
	lev.ExtraCalldepth = 2
	lev.SetBackend(backendLeveled)
	return lev
}

func (lg *ViLogsManage) vilLogErr(msg string, args ...interface{}) {
	logstr := fmt.Sprintf("%s message:%v", lg.LogConf.ProjectName, msg)
	lg.LogInfo.Errorf(logstr, args...)
}

func (lg *ViLogsManage) vilLogInfo(msg string, args ...interface{}) {
	logstr := fmt.Sprintf("%s message:%v", lg.LogConf.ProjectName, msg)
	lg.LogInfo.Infof(logstr, args...)
}

func (lg *ViLogsManage) vilLogwaring(msg string, args ...interface{}) {
	logstr := fmt.Sprintf("%s message:%v", lg.LogConf.ProjectName, msg)
	lg.LogInfo.Warningf(logstr, args...)
}

func (lg *ViLogsManage) vilLogDebug(msg string, args ...interface{}) {
	logstr := fmt.Sprintf("%s message:%v", lg.LogConf.ProjectName, msg)
	lg.LogInfo.Debugf(logstr, args...)
}

func LOGE(msg string, args ...interface{}) {
	if vilogs == nil {
		log.Println("error : Please call the function InitLogs before you call LOGE")
		panic(0)
	}
	vilogs.vilLogErr(msg, args...)
}

func LOGI(msg string, args ...interface{}) {
	if vilogs == nil {
		log.Println("error : Please call the function InitLogs before you call LOGI")
		panic(0)
	}
	vilogs.vilLogInfo(msg, args...)
}

func LOGW(msg string, args ...interface{}) {
	if vilogs == nil {
		log.Println("error : Please call the function InitLogs before you call LOGW")
		panic(0)
	}
	vilogs.vilLogwaring(msg, args...)
}

func LOGD(msg string, args ...interface{}) {
	if vilogs == nil {
		log.Println("error : Please call the function InitLogs before you call LOGD")
		panic(0)
	}
	vilogs.vilLogDebug(msg, args...)
}
