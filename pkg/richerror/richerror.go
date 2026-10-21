package richerror


type Kind int

const (
	KindInvalid Kind = iota+1
	KindForbidden 
	KindNotFound
	KindUnexpected
)


type RichError struct {
	operation	string
	wrappedError error
	message string
	kind Kind
	meta map[string]interface{}
}

func New(op string) RichError{
	return RichError{operation: op}
}

func (r RichError) WithOP(op string) RichError {
	r.operation = op
	return r
}

func (r RichError) WithErr(err error) RichError {
	r.wrappedError = err
	return r
}

func (r RichError) WithMessage(message string) RichError {
	r.message = message
	return r
}

func (r RichError) WithKind(Kind Kind) RichError {
	r.kind = Kind
	return r
}

func (r RichError) WithMeta(meta map[string]interface{}) RichError {
	r.meta = meta
	return r
}

func (r RichError) Error() string{
	return r.message
}

func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}
	
	re,ok := r.wrappedError.(RichError)
	if !ok{
		return 0
	}

	return re.Kind()
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}
	
	re,ok := r.wrappedError.(RichError)
	if !ok{
		return r.wrappedError.Error()
	}

	return re.Message()
}