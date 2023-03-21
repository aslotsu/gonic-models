package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gonic-models/database"
	"gonic-models/models"
	"log"
	"strconv"
	"time"
)

var client = database.Connect()

func addCountryHandlerFunc(c *gin.Context) {
	var country models.Country
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err := c.BindJSON(&country); err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}
	country.ID = primitive.NewObjectID()
	countriesColl := client.Database("world").Collection("countries")
	count, err := countriesColl.CountDocuments(ctx, bson.M{"name": country.Name})
	if count > 0 {
		c.JSON(402, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		c.JSON(400, gin.H{"err": err.Error()})
		return
	}

	result, err := countriesColl.InsertOne(ctx, country)
	if err != nil {
		c.JSON(404, gin.H{"err": err.Error()})
	}
	//
	c.JSON(200, result.InsertedID)
}

func getCountriesHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var countries []models.Country
	countriesColl := client.Database("world").Collection("countries")
	cursor, err := countriesColl.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Could not perform find operation on whole collection", err)
	}
	if err := cursor.All(ctx, &countries); err != nil {
		log.Println("Could not traverse collection", err)
		return
	}

	c.JSON(200, countries)
}

func geCountriesSortedHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var order = c.DefaultQuery("order", "1")
	sortOrder, err := strconv.Atoi(order)
	if err != nil {
		log.Println("Could not change string value to int", err)
		return
	}
	var countries []models.Country
	countriesColl := client.Database("world").Collection("countries")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"name", sortOrder}})
	cursor, err := countriesColl.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Println("Cannot perform find query", err)
		return
	}
	if err := cursor.All(ctx, &countries); err != nil {
		log.Println("Could not traverse collection", err)
		return
	}
	c.JSON(200, countries)

}

func getCountryHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var country models.Country
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Could not return objectID from hex")
	}
	countriesColl := client.Database("world").Collection("countries")
	if err := countriesColl.FindOne(ctx, bson.D{{"_id", objectId}}).Decode(&country); err != nil {
		log.Println("Could not return result of query", err)
		return
	}
	c.JSON(200, country)

}

func deleteCountryHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Could not return objectID from hex")
		return
	}
	countriesColl := client.Database("world").Collection("countries")
	result, err := countriesColl.DeleteOne(ctx, bson.D{{"_id", objectId}})
	if err != nil {
		log.Println("Unable to perform delete operation", err)
		return
	}
	c.JSON(200, result)
}

func updateCountryHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var country models.Country
	if err := c.BindJSON(&country); err != nil {
		log.Println("Could not bing request body")
		return
	}
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Could not extract objectID from hex")
	}

	countriesColl := client.Database("world").Collection("countries")
	updateResult, err := countriesColl.UpdateByID(ctx,
		objectID, bson.D{{"$set", bson.D{{"name", country.Name},
			{"age", country.Age},
			{"population", country.Population},
			{"continent", country.Continent},
		}}})

	if err != nil {
		log.Println("Could not update selected document", err)
		return
	}
	c.JSON(200, updateResult.ModifiedCount)

}

func replaceCountryHandlerFunc(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Could not extract objectID from hex", err)
		return
	}
	var country models.Country
	if err := c.BindJSON(&country); err != nil {
		log.Println("Could not bind json request body", err)
		return
	}

	countriesColl := client.Database("world").Collection("countries")
	updateResult, err := countriesColl.ReplaceOne(ctx,
		bson.D{{"_id", objectID}}, bson.D{{"name", country.Name},
			{"age", country.Age}, {"population", country.Population},
			{"continent", country.Continent}})
	if err != nil {
		log.Println("Could not perform replace operation", err)
		return
	}
	c.JSON(200, updateResult.ModifiedCount)

}

// AddCountry /*Add a new country to db*/
func AddCountry() gin.HandlerFunc {
	return addCountryHandlerFunc
}

// GetCountries /* Get All Countries*/
func GetCountries() gin.HandlerFunc {
	return getCountriesHandlerFunc
}

// GetCountriesSorted /*Get countries sorted */
func GetCountriesSorted() gin.HandlerFunc {
	return geCountriesSortedHandlerFunc
}

// GetCountry /* Get  Countries By ID*/
func GetCountry() gin.HandlerFunc {
	return getCountryHandlerFunc
}

// DeleteCountry /* Delete Country By ID*/
func DeleteCountry() gin.HandlerFunc {
	return deleteCountryHandlerFunc
}

// UpdateCountry /* Update Country By ID*/
func UpdateCountry() gin.HandlerFunc {
	return updateCountryHandlerFunc
}

// ReplaceCountry /*Replace country document*/
func ReplaceCountry() gin.HandlerFunc {
	return replaceCountryHandlerFunc
}
