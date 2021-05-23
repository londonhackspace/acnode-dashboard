package acnode

// Tracks Octoprint status
type PrinterStatus struct {
	MqttConneced bool

	// each status field has a value and a timestamp - this allows us to figure out which
	// topic has the newest information on it
	OctoprintConnected bool
	OctoprintConnectedTimestamp uint64

	PrinterStatus string
	PrinterStatusTimestamp uint64

	FirmwareVersion string
	FirmwareVersionTimestamp uint64

	ZHeight float32
	ZHeightTimestamp uint64

	PiUndervoltage bool
	PiOverheat bool
	PiThrottleTimestamp uint64
	
	HotendTemperature float32
	HotendTemperatureTimestamp uint64

	BedTemperature float32
	BedTemperatureTimestamp uint64
}

func GetDefaultPrinterStatus() *PrinterStatus {
	return &PrinterStatus{
		MqttConneced: false,
		OctoprintConnected: false,
		OctoprintConnectedTimestamp: 0,
		PrinterStatus: "",
		PrinterStatusTimestamp: 0,
		FirmwareVersion: "",
		FirmwareVersionTimestamp: 0,
		ZHeight: 0.0,
		ZHeightTimestamp: 0,
		PiUndervoltage: false,
		PiOverheat: false,
		PiThrottleTimestamp: 0,
		HotendTemperature: 0.0,
		HotendTemperatureTimestamp: 0,
		BedTemperature: 0.0,
		BedTemperatureTimestamp: 0,
	}
}