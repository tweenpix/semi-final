package lib

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

type City struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Region     string `json:"region"`
	District   string `json:"district"`
	Population int    `json:"population"`
	Foundation int    `json:"foundation"`
}

type CityService struct {
	Cities []City
}

func NewCityService(filename string) (*CityService, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	var cities []City
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		id, _ := strconv.Atoi(record[0])
		population, _ := strconv.Atoi(record[4])
		foundation, _ := strconv.Atoi(record[5])

		city := City{
			ID:         id,
			Name:       record[1],
			Region:     record[2],
			District:   record[3],
			Population: population,
			Foundation: foundation,
		}

		cities = append(cities, city)
	}

	return &CityService{
		Cities: cities,
	}, nil
}

func (s *CityService) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ','

	for _, city := range s.Cities {
		record := []string{
			strconv.Itoa(city.ID),
			city.Name,
			city.Region,
			city.District,
			strconv.Itoa(city.Population),
			strconv.Itoa(city.Foundation),
		}

		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	writer.Flush()

	return nil
}

func (s *CityService) GetCityByID(w http.ResponseWriter, r *http.Request) {

	// Извлекаем ID из URL-адреса
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, city := range s.Cities {
		if city.ID == id {
			fmt.Fprintf(w, `По ID %d найден город %v`, id, city.Name)
		}
	}

}

func (s *CityService) AddCity(w http.ResponseWriter, r *http.Request) {
	var city City

	err := json.NewDecoder(r.Body).Decode(&city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	s.Cities = append(s.Cities, city)

}

func (s *CityService) DeleteCityByID(w http.ResponseWriter, r *http.Request) {

	// Извлекаем ID из URL-адреса
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, city := range s.Cities {
		if city.ID == id {
			s.Cities = append(s.Cities[:i], s.Cities[i+1:]...)
			fmt.Fprintf(w, `Город %v удален`, city.Name)
			break
		}
	}
}

func (s *CityService) UpdatePopulationByID(w http.ResponseWriter, r *http.Request) {

	var populationRequest struct {
		Population int `json:"population"`
	}

	err := json.NewDecoder(r.Body).Decode(&populationRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	// Извлекаем ID из URL-адреса
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, city := range s.Cities {
		if city.ID == id {
			s.Cities[i].Population = populationRequest.Population
			fmt.Fprintf(w, `Численность населения в городе %v обновлена`, city.Name)
			break
		}
	}
}

func (s *CityService) GetCitiesByRegion(w http.ResponseWriter, r *http.Request) {

	var cities []City

	var regionRequest struct {
		Region string `json:"region"`
	}

	err := json.NewDecoder(r.Body).Decode(&regionRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	for _, city := range s.Cities {
		if city.Region == regionRequest.Region {
			cities = append(cities, city)
			fmt.Fprintf(w, "Город в %s: %s\n", regionRequest.Region, city.Name)
		}
	}

}

func (s *CityService) GetCitiesByDistrict(w http.ResponseWriter, r *http.Request) {
	var cities []City

	var Request struct {
		District string `json:"district"`
	}

	err := json.NewDecoder(r.Body).Decode(&Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	for _, city := range s.Cities {
		if city.District == Request.District {
			cities = append(cities, city)
			fmt.Fprintf(w, "Город в %s: %s\n", Request.District, city.Name)
		}
	}
}

func (s *CityService) GetCitiesByPopulationRange(w http.ResponseWriter, r *http.Request) {
	var cities []City
	var Request struct {
		PopulationMin int `json:"pmin"`
		PopulationMax int `json:"pmax"`
	}

	err := json.NewDecoder(r.Body).Decode(&Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	// fmt.Println(Request)

	for _, city := range s.Cities {
		if city.Population >= Request.PopulationMin && city.Population <= Request.PopulationMax {
			cities = append(cities, city)
			fmt.Fprintf(w, "Город с численностью населения %d - %d: %s\n", Request.PopulationMin, Request.PopulationMax, city.Name)
		}
	}

}

func (s *CityService) GetCitiesByFoundationRange(w http.ResponseWriter, r *http.Request) {
	var cities []City
	var Request struct {
		FoundationMin int `json:"fmin"`
		FoundationMax int `json:"fmax"`
	}

	err := json.NewDecoder(r.Body).Decode(&Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	for _, city := range s.Cities {
		if city.Foundation >= Request.FoundationMin && city.Foundation <= Request.FoundationMax {
			cities = append(cities, city)
			fmt.Fprintf(w, "Город с годом основания %d - %d: %s\n", Request.FoundationMin, Request.FoundationMax, city.Name)

		}
	}

}
