package spacefile

type ResourceDef struct {
	Resource string      `hcl:",key"`
	Quantity interface{} `hcl:"quantity"`
}
