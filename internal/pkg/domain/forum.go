package domain

//easyjson:json
type Forum struct {
	Id      int
	Title   string
	Slug    string
	User    string
	Posts   int
	Threads int
}
