package api

import (
	"net/http"
	"regexp"
	"time"
	utils "vercelgin/api/_pkg"

	"github.com/gin-gonic/gin"
)

var (
	app    *gin.Engine
	re     = regexp.MustCompile(`(?i)[A-Za-z]\d\d\d\d\d\d\d\d\d`) // Make it global so we don't recompile every time
	client = &http.Client{                                        // We can use the same client again and again
		Timeout: 5 * time.Second,
	}
)

// CREATE ENDPOIND

func myRoute(r *gin.RouterGroup) {
	r.GET("/v1", func(c *gin.Context) {
		// Allow fetch api
		c.Header("Access-Control-Allow-Origin", "*")
		//Get CNE
		cne := c.Query("CNE")

		// Validate CNE format
		isValid := re.MatchString(cne)
		// Non valid CNE then return Json error
		if !isValid {
			c.JSON(404, gin.H{"error": "Invalid CNE Format"})
			return
		}

		// Is a Valid CNE, then make a request
		req := utils.MakeRequest(cne)

		// Send the request using the global client
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Cannot get Data please try later"})
			return
		}
		defer resp.Body.Close()

		// We can parse our resp.body now
		status, data := utils.ParseBody(resp.Body)
		c.JSON(status, data)
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
