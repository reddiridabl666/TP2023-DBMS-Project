package delivery

//easyjson:json
type ErrorMessage struct {
	Message string
}

var MsgBadJSON = ErrorMessage{
	Message: "bad JSON got",
}

var MsgUserNotFound = ErrorMessage{
	Message: "no such user",
}

var MsgUserExists = ErrorMessage{
	Message: "such user already exists",
}

var MsgInternalError = ErrorMessage{
	Message: "internal error",
}
