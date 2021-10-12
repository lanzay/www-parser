package auto_ru

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/lanzay/www-parser/models"
	"github.com/lanzay/www-parser/tools"
	"golang.org/x/net/html"
	"log"
	"strings"
)

const (
	ENDPOINT = "https://auto.ru/"
)

func GetBrands() []models.Item {

	method := "catalog/cars/"
	code, body, err := tools.GetBody(ENDPOINT + method)
	if code != 200 || err != nil {
		log.Panicln("[E]", code, err)
	}

	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	doc := goquery.NewDocumentFromNode(root)

	var items []models.Item
	doc.Find("div.search-form-v2-list__text-item>a.i-bem").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		item := models.Item{
			Name: s.Text(),
			URL:  href,
		}
		//log.Println("[D]", item)
		items = append(items, item)
	})

	return items
}

func GetModelsByMark(mark, uri string) []models.Item {

	if uri == "" {
		method := "catalog/cars/" + strings.ToLower(mark) + "/"
		uri = ENDPOINT + method
	}

	code, body, err := tools.GetBody(uri)
	if code != 200 || err != nil {
		log.Panicln("[E]", code, err)
	}

	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	doc := goquery.NewDocumentFromNode(root)

	var items []models.Item
	doc.Find("div.search-form-v2-list__text-column>div.search-form-v2-list__text-item>a.i-bem").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		item := models.Item{
			Name: s.Text(),
			URL:  href,
		}
		//log.Println("[D]", item)
		items = append(items, item)
	})

	return items
}

func GetGenerations(mark, model, uri string) []models.Generation {

	if uri == "" {
		model = strings.ReplaceAll(model, "-", "_")
		model = strings.ReplaceAll(model, " ", "_")

		method := "catalog/cars/" + strings.ToLower(mark) + "/" + strings.ToLower(model) + "/"
		uri = ENDPOINT + method
	}

	code, body, err := tools.GetBody(uri)
	if code != 200 || err != nil {
		log.Println("[E]", code, err)
		return nil
	}

	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	doc := goquery.NewDocumentFromNode(root)

	var items []models.Generation
	//doc.Find("dl.catalog-all-text-list>dd>div>div>a").Each(func(i int, s *goquery.Selection) {
	doc.Find("dl.catalog-all-text-list").Each(func(i int, s *goquery.Selection) {
		dt := s.Find("dt")
		dd := s.Find("dd")
		for i := range dd.Nodes {
			print(i, dt.Nodes[i])
			item := models.Generation{
				Name:  goquery.NewDocumentFromNode(dt.Nodes[i]).Find("div").Text(),
				Years: dt.Nodes[i].FirstChild.Data,
				URL:   goquery.NewDocumentFromNode(dd.Nodes[i]).Find("a").AttrOr("href",""),
			}
			items = append(items, item)
		}
	})

	return items
}

func GetSpecificationsByURL(u string) (*models.Car, []models.Modification) {

	method := "specifications/"
	uri := u
	if !strings.Contains(u, method) {
		uri = u + method
	}

	code, body, err := tools.GetBody(uri)
	if code != 200 || err != nil {
		log.Panicln("[E]", code, err)
	}

	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	doc := goquery.NewDocumentFromNode(root)

	// Модификации
	var modifications []models.Modification
	{
		fuel := ""
		complectationName := ""
		// .catalog__sidebar>.catalog-table>.catalog-table__row - Модификации
		doc.Find(".catalog-table__row").Each(func(i int, s *goquery.Selection) {
			if tmp := s.Find(".catalog-table__group-title").Text(); tmp != "" {
				fuel = tmp
				return
			}
			if tmp := s.Find(".catalog-table__label-title").Text(); tmp != "" {
				complectationName = tmp
				if complectationName == "–" {
					complectationName = ""
				}
				return
			}
			a := s.Find("a").First()
			href, _ := a.Attr("href")
			item := models.Modification{
				Name:              a.Text(),
				ComplectationName: complectationName,
				Fuel:              fuel,
				Power:             s.Find(".catalog-table__cell_alias_power").Text(),
				Gear:              s.Find(".catalog-table__cell_alias_gear").Text(),
				URL:               href,
			}
			modifications = append(modifications, item)
		})
	}

	// Комплектация/Модификация
	complectation := models.Complectation{
		URL: uri,
	}
	{
		// .catalog__content>.catalog__details-main - Модификация
		doc.Find(".catalog__content>.catalog__details-main").Each(func(i int, s *goquery.Selection) {
			complectation.Name = s.Find("h2").Text()
			dt := s.Find(".list-values>dt")
			dd := s.Find(".list-values>dd")
			for i := range dd.Nodes {
				characteristic := models.Specification{
					Head:  "Основные",
					Name:  dt.Nodes[i].FirstChild.Data,
					Value: dd.Nodes[i].FirstChild.Data,
				}
				complectation.Specification = append(complectation.Specification, characteristic)
			}
		})
	}

	{
		// .catalog__content>.clearfix>.catalog__column>.catalog__details-group - Характеристики
		doc.Find(".catalog__content>.clearfix>.catalog__column>.catalog__details-group").Each(func(i int, s *goquery.Selection) {
			title := s.Find("h3").Text()
			dt := s.Find(".list-values>dt")
			dd := s.Find(".list-values>dd")
			for i := range dd.Nodes {
				characteristic := models.Specification{
					Head:  title,
					Name:  dt.Nodes[i].FirstChild.Data,
					Value: dd.Nodes[i].FirstChild.Data,
				}
				complectation.Specification = append(complectation.Specification, characteristic)
			}
		})
	}

	labels := strings.Split(uri, "/")
	car := &models.Car{
		Brand: labels[5],
		Model: labels[6],
		//Modification:  modification,
		Complectation: complectation,
	}
	return car, modifications
}
