package pokemon

type pokemonResponse struct {
	Name    string                 `json:"name"`
	Id      int                    `json:"id"`
	Sprites pokemonSpritesResponse `json:"sprites"`
}

type pokemonSpritesResponse struct {
	Default     string `json:"front_default"`
	BackDefault string `json:"back_default"`
	FrontShiny  string `json:"front_shiny"`
}

type pokemonListResponse struct {
	Results []pokemonListResultResponse `json:"results"`
}

type pokemonListResultResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
