package spaces

func (c *spacesClient) Delete(spaceID string) (error) {
	var target interface{}
	return c.client.Delete("/spaces/" + spaceID, target)
}
