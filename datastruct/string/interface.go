package string

type String interface {
	Set(string)
	Get() string
	Append(string)
	Len() int
}
