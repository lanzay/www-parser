package main

import (
	"encoding/json"
	"github.com/lanzay/www-parser/target/auto_ru"
	"log"
	"os"
	"strings"
)

func main() {

	skip := true
	marks := auto_ru.GetBrands()
	for _, mark := range marks {
		if strings.EqualFold(mark.Name,"Mazda") {
			skip = false
		}
		if skip {
			continue
		}

		fn, _ := os.Create("auto_ru_" + mark.Name + ".json")
		defer fn.Close()

		models := auto_ru.GetModelsByMark(mark.Name,mark.URL)
		for _, model := range models {
			gens := auto_ru.GetGenerations(mark.Name, model.Name, model.URL)
			for _, gen := range gens {
				_, complectations := auto_ru.GetSpecificationsByURL(gen.URL)
				for _, complectation := range complectations {
					car, _ := auto_ru.GetSpecificationsByURL(complectation.URL)
					car.Brand = mark.Name
					car.Model = model.Name
					car.Generation = gen
					body, _ := json.Marshal(car)
					fn.Write(body)
					fn.Write([]byte("\n"))
					fn.Sync()

					log.Println("[D]", car.Brand, car.Model, car.Complectation.Name, car.Generation)
				}
			}
		}
		fn.Close()
	}
}
