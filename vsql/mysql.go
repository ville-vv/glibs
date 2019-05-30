package vsql

const (
	MaxIdleConnsDefautl = 100
	MaxOpenConnsDefault = 1000
)

type MySqlCnf struct {
	Version      string   `json:"version"`
	UserName     string   `json:"user_name"`
	Address      string   `json:"host"`
	Password     string   `json:"password"`
	Default      string   `json:"default"`
	MaxIdleConns int      `json:"max_idle_conns"`
	MaxOpenConns int      `json:"max_open_conns"`
	Databases    []string `json:"databases"`
}
