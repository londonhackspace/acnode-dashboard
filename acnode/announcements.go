package acnode

type StatusMem struct {
	HeapFree int `json:"heap_free",omitempty`
	HeapUsed int `json:"heap_used",omitempty`
}

type Status struct {
	Message string `json:"Message",omitempty`
	Mem StatusMem `json:"mem",omitempty`
}

type Announcement struct {
	Type string `json:"Type"`
	Card string `json:"Card",omitempty`
	Granted int `json:"Granted",omitempty`
}