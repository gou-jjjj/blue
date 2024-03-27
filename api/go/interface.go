package blue

type CLi interface {
	//system
	Version() (string, error)
	Select(...string) (string, error)
	Exit() (string, error)
	Ping() (string, error)
	Help(string) (string, error)

	//number
	Nset(string, string) (string, error)
	Nget(string) (string, error)
	Incr(string) (string, error)

	//db
	Len(string) (string, error)
	Kvs() (string, error)
	Del(string) (string, error)
	Expire(string, string) (string, error)

	//list
	Llen(string) (string, error)
	Lget(string) (string, error)
	Lset(string) (string, error)
	Lpop(string) (string, error)
	Lpush(string, string) (string, error)
	Rpop(string) (string, error)
	Rpush(string, string) (string, error)

	// Set
	Sadd(string, string) (string, error)
	Spop(string) (string, error)
	Sin(string, string) (string, error)
	Sdel(string, string) (string, error)
	Sget(string) (string, error)

	// String
	Get(string) (string, error)
	Set(string, string) (string, error)
}
