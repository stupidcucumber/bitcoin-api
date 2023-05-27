package main

import (
	"net/http"
	"os"

	"log"

	"api/bitcoin-api/controlers"
	"api/bitcoin-api/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getPrice(c *gin.Context) {
	answer := make(map[string]float64)

	price, err := controlers.GetPrice()

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	answer["rate"] = price

	c.IndentedJSON(http.StatusOK, answer)
}

func postSubscribe(c *gin.Context) {
	var newEmail models.Email

	if err := c.BindJSON(&newEmail); err != nil {
		c.IndentedJSON(http.StatusConflict, newEmail)
		return
	}

	if err := controlers.AddEmail(newEmail); err == nil {
		c.IndentedJSON(http.StatusOK, newEmail)
	} else {
		c.IndentedJSON(http.StatusConflict, err.Error())
	}
}

func postSendEmails(c *gin.Context) {
	price, _ := controlers.GetPrice()

	controlers.SendEmail(price)

	c.IndentedJSON(http.StatusOK, "Emails had been sent")
}

func main() {
	router := gin.Default()

	router.GET("/rate", getPrice)
	router.POST("/subscribe", postSubscribe)
	router.POST("/sendEmails", postSendEmails)

	router.Run(":" + os.Getenv("APP_PORT"))
}
