package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

type car struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Model  string `json:"model"`
	Run    int    `json:"run"`
	Owners byte   `json:"owners"`
}

var cars []car

func loadCarsFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cars)
}

// saveCarsToFile сохраняет машины в  JSON файл.
func saveCarsToFile(filename string) error {
	data, err := json.MarshalIndent(cars, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// postCars добавляет машину из JSON, полученного в теле запроса.
func postCars(c *gin.Context) {
	var newCar car

	// Вызываем BindJSON, чтобы привязать полученный JSON к
	// newCar.
	if err := c.BindJSON(&newCar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Добавляем новую машину в срез.
	cars = append(cars, newCar)

	// Сохраняет машины в файл
	if err := saveCarsToFile("cars.json"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newCar)
}

// getCarByID находит машину, значение id которой совпадает с параметром id
// , отправленным клиентом, и возвращает эту машину в качестве ответа.
func getCarByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range cars {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "car not found"})
}

// updateCarByID обновляет поля машины , заданного клиентом id.
func updateCarByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range cars {
		if a.ID == id {
			var updatedCar car
			if err := c.BindJSON(&updatedCar); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			cars[i] = updatedCar

			if err := saveCarsToFile("cars.json"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.IndentedJSON(http.StatusOK, cars[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "car not found"})
}

// patchCarByID частично обновляет поля машины с заданным ID.
func patchCarByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range cars {
		if a.ID == id {

			var patchCar car
			if err := c.BindJSON(&patchCar); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if patchCar.Name != "" {
				cars[i].Name = patchCar.Name
			}
			if patchCar.Model != "" {
				cars[i].Model = patchCar.Model
			}
			if patchCar.Run != 0 {
				cars[i].Run = patchCar.Run
			}
			if patchCar.Owners != 0 {
				cars[i].Owners = patchCar.Owners
			}

			if err := saveCarsToFile("cars.json"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.IndentedJSON(http.StatusOK, cars[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "car not found"})
}

// deleteCarByID удаляет машину с заданным ID.
func deleteCarByID(c *gin.Context) {

	id := c.Param("id")
	for i, a := range cars {
		if a.ID == id {

			cars = append(cars[:i], cars[i+1:]...)

			if err := saveCarsToFile("cars.json"); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.Status(http.StatusNoContent) // Return 204 No Content
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "car not found"})
}

func main() {
	// Загружаем машины из файла при запуске.
	if err := loadCarsFromFile("cars.json"); err != nil {
		// Если файл не существует, создаем новый с машинами по умолчанию
		if os.IsNotExist(err) {
			fmt.Println("File not found, creating a new one with default albums.")
			cars = []car{
				{ID: "1", Name: "Toyota", Model: "Rav 4", Run: 100000, Owners: 3},
				{ID: "2", Name: "BMW", Model: "3-Series", Run: 50000, Owners: 1},
				{ID: "3", Name: "Haval", Model: "M6", Run: 100000, Owners: 2},
			}
			if err := saveCarsToFile("cars.json"); err != nil {
				fmt.Println("Error saving cars to file:", err)
			}
		} else {
			fmt.Println("Error loading cars from file:", err)
		}
	}
	router := gin.Default()
	router.GET("/cars", getCars)
	router.POST("/cars", postCars)
	router.GET("/cars/:id", getCarByID)
	router.PUT("/cars/:id", updateCarByID)
	router.PATCH("/cars/:id", patchCarByID)
	router.DELETE("/cars/:id", deleteCarByID)
	router.Run("localhost:8080")
}

// getCars выдает список всех машин в формате JSON.
func getCars(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cars)
}
