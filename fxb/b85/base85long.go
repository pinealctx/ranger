package b85

type B85LongConverter interface {
	AsString(int64) string
	Parse(string) (int64, error)
}
