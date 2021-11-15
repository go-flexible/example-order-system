package domain

type PublishingError struct{ error }

type DatabaseQueryError struct {
	Inner error
	Stmt  string
}

func (e DatabaseQueryError) Error() string {
	return e.Inner.Error()
}
