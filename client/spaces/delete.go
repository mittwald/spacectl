package spaces

type Empty struct {}

func (c *spacesClient) Delete(spaceID string) (error) {
	var target Empty
	return c.client.Delete("/spaces/" + spaceID, &target)
}
