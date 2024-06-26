package serve

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go-salesforce-crosscloud/server"
	"log"
	"time"
)

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Salesforce crosscloud server",
		Long:  "Salesforce crosscloud server",
		Run:   runServe,
	}

	return cmd
}

func runServe(cmd *cobra.Command, args []string) {
	router := gin.New()
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "x-access-token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	})

	router.Use(corsMiddleware)
	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "Not found", "message": "Page not found"})
	})

	router.POST("/", server.HandleForm)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
}
