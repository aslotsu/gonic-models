package main

import (
	"github.com/gin-gonic/gin"
	"gonic-models/routes"
	"log"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "You were successful in getting the Api 🇬🇭"})

	})

	router.POST("/new", routes.AddCountry())
	router.GET("/all", routes.GetCountries())
	router.GET("/all/sorted", routes.GetCountriesSorted())
	router.GET("/all/:id", routes.GetCountry())
	router.DELETE("/del/:id", routes.DeleteCountry())
	router.PUT("/repl/:id", routes.ReplaceCountry())
	router.PATCH("/upd/:id", routes.UpdateCountry())
	if err := router.Run(":8083"); err != nil {
		log.Println("Could not start router at said port", err)
		log.Fatal()
	}
}
