package controllers

type ClientDetails struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

func (c *Controller) updateClientDetails(token string, details ClientDetails) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ClientCache[token] = details
}

// Function to retrieve user details from the map
func (c *Controller) getClientDetails(token string) (ClientDetails, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	details, exists := c.ClientCache[token]
	return details, exists
}
