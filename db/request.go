package db

type Request struct {
	Query      string
	Args       []any
	ResultChan chan Result
}

type Result struct {
	Data  any
	Error error
}
