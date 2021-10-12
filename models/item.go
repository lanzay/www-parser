package models

type Item struct {
	Name string
	URL  string
}

type Modification struct {
	Name              string
	ComplectationName string
	Fuel              string
	Power             string
	Gear              string
	URL               string
}

type (
	Complectation struct {
		Name          string `json:"name"`
		Specification []Specification `json:"specification"`
		URL           string `json:"url"`
	}
	Specification struct {
		Head  string `json:"head"`
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	Generation struct {
		Name  string `json:"name"`
		Years string `json:"years"`
		URL   string `json:"url"`
	}
)

type Car struct {
	Brand string `json:"brand"`
	Model string `json:"model"`
	Generation    Generation `json:"generation"`
	Complectation Complectation `json:"complectation"`
	//Modification  Modification
}
