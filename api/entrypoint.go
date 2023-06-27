package api

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
	utils "vercelgin/api/_pkg"

	"github.com/gin-gonic/gin"
)

// Data that we try to grab
type Data struct {
	CIN         string `json:"cin"`
	CNE         string `json:"cne"`
	ARLastName  string `json:"last_name_ar"`
	FRLastName  string `json:"last_name_fr"`
	ARFirstName string `json:"first_name_ar"`
	FRFirstName string `json:"first_name_fr"`
	BirthDate   string `json:"birth_date"`
}

var (
	app *gin.Engine
)

// CREATE ENDPOIND

func myRoute(r *gin.RouterGroup) {
	r.GET("/v1", func(c *gin.Context) {
		// Initialise regex
		re := regexp.MustCompile(`(?i)[A-Za-z]\d\d\d\d\d\d\d\d\d`)

		//Get CNE
		cne := c.Query("CNE")

		// Validate this CNE
		isValid := re.MatchString(cne)
		// Non valid CNE then return Json error
		if !isValid {
			c.JSON(404, gin.H{"error": "Invalid CNE Format"})
			return
		}

		// Is a Valid CNE, then make a request
		req := utils.MakeRequest(cne)

		// make a client with short timeout because in serverless time cost money
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Cannot get Data please try later"})
			return
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Cannot turn body to bytes"})
			return
		}

		// convert the byte slice to a string
		c.String(http.StatusOK, string(bytes))
	})

}

func init() {
	app = gin.New()
	r := app.Group("/api")
	myRoute(r)

}

// ADD THIS SCRIPT
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
