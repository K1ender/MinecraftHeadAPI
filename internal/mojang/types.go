package mojang

type MojangSessionResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type texturesPayload struct {
	Textures struct {
		SKIN struct {
			URL      string `json:"url"`
			Metadata struct {
				Model string `json:"model,omitempty"`
			} `json:"metadata,omitempty"`
		} `json:"SKIN"`
	} `json:"textures"`
}
