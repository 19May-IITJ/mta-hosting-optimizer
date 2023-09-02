package ipconfig

var _ IPListI = (*IPConfigs)(nil)

type IPConfigs struct {
	ipConfigList []*IPConfigData
}

// Factory Method returns New IPConfigs object
func NewIPConfigList() *IPConfigs {
	return &IPConfigs{
		ipConfigList: make([]*IPConfigData, 0),
	}
}

// Getter of the IPList
func (c *IPConfigs) GetIPList() *IPConfigs {
	return c
}

// Setter of the IPList
func (c *IPConfigs) SetIPList(l []*IPConfigData) {
	c.ipConfigList = append(c.ipConfigList, l...)
}

// Checks whether IPList is Empty or not
func (c *IPConfigs) IsEmpty() bool {
	return len(c.GetIPList().ipConfigList) == 0
}

// Clears the IPList
func (c *IPConfigs) Clear() {
	c.ipConfigList = make([]*IPConfigData, 0)
}

// Returns the size of underlying IPlist
func (c *IPConfigs) Size() int {
	return len(c.GetIPList().ipConfigList)
}

// Getter of the underlying list of IP
func (c *IPConfigs) GetIPValues() []*IPConfigData {
	return c.ipConfigList
}
