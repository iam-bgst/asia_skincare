package swagger

type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Info struct {
	Description    string  `json:"description"`
	Title          string  `json:"title"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
}

type PathGET struct{
	Get struct {
		Description string   `json:"description"`
		Consumes    []string `json:"consumes"`
		Produces    []string `json:"produces"`
		Summary     string   `json:"summary"`
		Parameters  []struct {
			Type        string `json:"type"`
			Description string `json:"description"`
			Name        string `json:"name"`
			In          string `json:"in"`
			Required    bool   `json:"required"`
		} `json:"parameters"`
		Responses struct {
			Success struct {
				Description string `json:"description"`
				Schema      struct {
					Type string `json:"type"`
				} `json:"schema"`
			} `json:"200"`
			Warning struct {
				Description string `json:"description"`
				Schema      struct {
					Ref string `json:"$ref"`
				} `json:"schema"`
			} `json:"400"`
			Failed struct {
				Description string `json:"description"`
				Schema      struct {
					Ref string `json:"$ref"`
				} `json:"schema"`
			} `json:"404"`
		} `json:"get"`
	}
}
type SwaggerConfig struct {
	Schemes  []string `json:"schemes"`
	Swagger  string   `json:"swagger"`
	Info     Info     `json:"info"`
	Host     string   `json:"host"`
	BasePath string   `json:"basePath"`
	Paths    struct {
		TestapiGetStringByIntSomeID struct {
			Get struct {
				 `json:"responses"`
			} `json:"get"`
		} `json:"/testapi/get-string-by-int/{some_id}"`
		TestapiGetStructArrayByStringSomeID struct {
			Get struct {
				Description string   `json:"description"`
				Consumes    []string `json:"consumes"`
				Produces    []string `json:"produces"`
				Parameters  []struct {
					Type        string `json:"type"`
					Description string `json:"description"`
					Name        string `json:"name"`
					In          string `json:"in"`
					Required    bool   `json:"required"`
				} `json:"parameters"`
				Responses struct {
					Num200 struct {
						Description string `json:"description"`
						Schema      struct {
							Type string `json:"type"`
						} `json:"schema"`
					} `json:"200"`
					Num400 struct {
						Description string `json:"description"`
						Schema      struct {
							Ref string `json:"$ref"`
						} `json:"schema"`
					} `json:"400"`
					Num404 struct {
						Description string `json:"description"`
						Schema      struct {
							Ref string `json:"$ref"`
						} `json:"schema"`
					} `json:"404"`
				} `json:"responses"`
			} `json:"get"`
		} `json:"/testapi/get-struct-array-by-string/{some_id}"`
	} `json:"paths"`
	Definitions struct {
		WebAPIError struct {
			Type       string `json:"type"`
			Properties struct {
				ErrorCode struct {
					Type string `json:"type"`
				} `json:"errorCode"`
				ErrorMessage struct {
					Type string `json:"type"`
				} `json:"errorMessage"`
			} `json:"properties"`
		} `json:"web.APIError"`
	} `json:"definitions"`
}
