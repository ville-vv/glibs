package vredis

type RedisCnf struct {
	Address       string `json:"address"`
	Password      string `json:"password"`
	UserName      string `json:"user_name"`
	MaxIdles      int    `json:"max_idles"`
	MaxOpens      int    `json:"max_opens"`
	IdleTimeout   int    `json:"idle_timeout"`
	IsMaxConnWait bool   `json:"is_max_conn_wait"` // 达到最大连接数后，是否等待
	Db            int    `json:"db"`
}
