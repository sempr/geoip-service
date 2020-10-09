package main

//go:generate packr2
import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/oschwald/geoip2-golang"
)

func getDB() (*geoip2.Reader, error) {
	box := packr.New("GEO-DATA", "../../data/")
	bytes, err := box.Find("GeoLite2-Country.mmdb")
	if err != nil {
		return nil, err
	}
	db, err := geoip2.FromBytes(bytes)
	return db, err
}

func getOrigin(c *gin.Context) string {
	r := c.Request
	origin := r.Header.Get("X-Forwarded-For")
	if origin == "" {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		origin = ip
	} else {
		xffs := strings.SplitN(origin, ", ", 2)
		origin = xffs[0]
	}
	return origin
}

func main() {
	db, err := getDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/jsonp", func(c *gin.Context) {
		ipstr := getOrigin(c)
		ip := net.ParseIP(ipstr)
		record, err := db.Country(ip)
		if err != nil {
			c.JSONP(http.StatusInternalServerError, gin.H{})
		}
		data := gin.H{
			"ip":      ipstr,
			"country": record.Country.IsoCode,
		}
		c.JSONP(http.StatusOK, data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8008"
	}
	r.Run(":" + port)
}
