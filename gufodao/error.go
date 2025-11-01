package gufodao

func Error(code string) GufoError {
	if e, ok := Errors[code]; ok {
		return e
	}
	return GufoError{"99999", "Unknown Error", 500}
}

var Errors = map[string]GufoError{
	"00001": {"00001", "Unauthorized", 401},
	"00002": {"00002", "Invalid Session", 401},
	"00003": {"00003", "Bad Request", 400},
	"00004": {"00004", "Internal Error", 500},
}
