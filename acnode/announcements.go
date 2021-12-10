package acnode

type StatusMem struct {
	HeapFree int `json:"heap_free,omitempty"`
	HeapUsed int `json:"heap_used,omitempty"`
}

type Status struct {
	Type string `json:"Type"`

	Message string    `json:"Message,omitempty"`
	Mem     StatusMem `json:"mem,omitempty"`

	// START message
	FWVersion             string `json:"Version,omitempty"`
	GitHash               string `json:"Git,omitempty"`
	SettingsVersion       int    `json:"SettingsVersion,omitempty"`
	EEPROMSettingsVersion int    `json:"EEPROMSettingsVersion,omitempty"`
	ResetCause            string `json:"Cause,omitempty"`
}

type Announcement struct {
	Type string `json:"Type"`

	// RFID message
	Card    string `json:"Card,omitempty"`
	Granted int    `json:"Granted,omitempty"`

	// EXIT,WEDGED messages
	Message string `json:"Message,omitempty"`

	// EXIT message
	DoorbellAck bool `json:"doorbellack"`
}

type Bell struct {
	Message string `json:"Message,omitempty"`
}
