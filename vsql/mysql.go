package vsql

const (
	MaxIdleConnsDefautl = 100
	MaxOpenConnsDefault = 1000
)

type MySqlCnf struct {
	Version   string   `json:"version"`
	UserName  string   `json:"user_name"`
	Address   string   `json:"host"`
	Password  string   `json:"password"`
	Default   string   `json:"default"`
	MaxIdles  int      `json:"max_idles"`
	MaxOpens  int      `json:"max_opens"`
	Databases []string `json:"databases"`
}
