package ipconfig

type IPConfigs struct {
	IpConfigList []*IPConfig
}

func NewIPConfigList() *IPConfigs {
	return &IPConfigs{
		IpConfigList: make([]*IPConfig, 0),
	}
}

func (c *IPConfigs) GetIPList() *IPConfigs {
	return c
}
func (c *IPConfigs) SetIPList(l []*IPConfig) {
	c.IpConfigList = append(c.IpConfigList, l...)
}
func (c *IPConfigs) IsEmpty() bool {
	return len(c.GetIPList().IpConfigList) == 0
}
func (c *IPConfigs) Clear() {
	c.IpConfigList = make([]*IPConfig, 0)
}
