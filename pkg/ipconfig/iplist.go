package ipconfig

var _ IPListI = (*IPConfigs)(nil)

type IPConfigs struct {
	ipConfigList []*IPConfig
}

// Factory Method returns New IPConfigs object
func NewIPConfigList() *IPConfigs {
	return &IPConfigs{
		ipConfigList: make([]*IPConfig, 0),
	}
}

// Getter of the IPList
func (c *IPConfigs) GetIPList() *IPConfigs {
	return c
}

// Setter of the IPList
func (c *IPConfigs) SetIPList(l []*IPConfig) {
	c.ipConfigList = append(c.ipConfigList, l...)
}

// Checks whether IPList is Empty or not
func (c *IPConfigs) IsEmpty() bool {
	return len(c.GetIPList().ipConfigList) == 0
}

// Clears the IPList
func (c *IPConfigs) Clear() {
	c.ipConfigList = make([]*IPConfig, 0)
}

// Returns the size of underlying IPlist
func (c *IPConfigs) Size() int {
	return len(c.GetIPList().ipConfigList)
}

// Getter of the underlying list of IP
func (c *IPConfigs) GetIPValues() []*IPConfig {
	return c.ipConfigList
}
