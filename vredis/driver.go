package vredis

type Driver interface {
	Conn() error
	Close() error
}
