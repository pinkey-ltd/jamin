package usecase

type DeviceTemplateProtocol string

const (
	ProtocolMQTT DeviceTemplateProtocol = "mqtt"
	ProtocolHTTP DeviceTemplateProtocol = "http"
)

func (p DeviceTemplateProtocol) String() string {
	return string(p)
}

type DeviceTemplateEntity struct {
	ID       int64  `json:"id" sql:"auto_increment"`
	Name     string `json:"name" sql:"not_null"`
	Model    string // The model of the device, representing its specific version or type.
	Maker    string // The manufacturer or brand that produced the device.
	Types    string // Specifies the category or type of the device, such as 'tester', 'camera', or 'hyetometer'.
	Protocol DeviceTemplateProtocol
}

func (t *DeviceTemplateEntity) RestPath() string {
	return "device_template"
}

func (t *DeviceTemplateEntity) TableName() string {
	return "mqtt_device_template"
}

type DeviceEntity struct {
	ID         int64  `json:"id" sql:"auto_increment"`
	Name       string `json:"name" sql:"not_null"`
	Status     string
	Health     string
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	TemplateID int64
	ParentID   int64
	JcdID      int64
}

func (u *DeviceEntity) RestPath() string {
	return "device"
}

func (u *DeviceEntity) TableName() string {
	return "mqtt_device"
}
