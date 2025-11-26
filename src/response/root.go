package resp

type Root[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func OK[T any](code int, data T) Root[T] {
	return Root[T]{
		Code:    code,
		Message: "success",
		Data:    data,
	}
}

func Error(code int, message string) Root[struct{}] {
	return Root[struct{}]{
		Code:    code,
		Message: message,
		Data:    struct{}{},
	}
}
