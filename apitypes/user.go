package apitypes

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
}
