package payment

type SpaceResourcePreprovisioningInputItem struct {
	Quantity uint64 `json:"quantity"`
}

type SpaceResourcePreprovisioningInput struct {
	Storage *SpaceResourcePreprovisioningInputItem `json:"storage,omitempty"`
	Stages  *SpaceResourcePreprovisioningInputItem `json:"stages,omitempty"`
	Scaling *SpaceResourcePreprovisioningInputItem `json:"scaling,omitempty"`
}
