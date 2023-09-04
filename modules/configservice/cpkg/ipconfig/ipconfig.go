package ipconfig

type IPConfigData struct {
	Hostname    string `json:"hostname"`
	IPAddresses string `json:"ipAddresses"`
	Status      bool   `json:"status"`
}

// type IPConfigData struct {
// 	Hostname    string `json:"hostname" validate:"required,string"`
// 	IPAddresses string `json:"ipAddresses" validate:"required,string"`
// 	Status      bool   `json:"status" validate:"required,boolean"`
// }
// "github.com/go-playground/validator/v10"
