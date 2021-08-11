package domain

type PublishingError struct{ error }

type DatabaseQueryError struct {
	Stmt  string
	Inner error
}

func (e DatabaseQueryError) Error() string {
	return e.Inner.Error()
}
