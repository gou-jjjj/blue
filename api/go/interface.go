package blue

type CLi interface {
	Version() (string, error)
	Select(...string) (string, error)
	Del(string) (string, error)
	Nset(string, string) (string, error)
	Get(string) (string, error)
	Set(string, string) (string, error)
	Len(string) (string, error)
	Kvs() (string, error)
	Nget(string) (string, error)
}
