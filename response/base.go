package response

type Body interface {
	Encode() ([]byte, error)
}
