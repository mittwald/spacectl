package spacefile

type OptionDef struct {
	Option string      `hcl:",key"`
	Value  interface{} `hcl:"value"`
}
