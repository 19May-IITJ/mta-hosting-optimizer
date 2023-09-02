package ipconfig

type IPConfigData struct {
	Hostname    string `json:"hostname"`
	IPAddresses string `json:"ipAddresses"`
	Status      bool   `json:"status"`
}
