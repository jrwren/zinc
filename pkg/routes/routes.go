package routes

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prabhatsharma/zinc/pkg/auth"
	"github.com/prabhatsharma/zinc/pkg/handlers"
	v1 "github.com/prabhatsharma/zinc/pkg/meta/v1"
	"github.com/rakyll/statik/fs"

	_ "github.com/prabhatsharma/zinc/statik"
)

// SetRoutes sets up all gi HTTP API endpoints that can be called by front end
func SetRoutes(r *gin.Engine) {

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "authorization", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// meta service - healthz
	r.GET("/healthz", v1.GetHealthz)
	r.GET("/", v1.GUI)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	r.StaticFS("/ui", statikFS)

	r.POST("/api/login", handlers.ValidateCredentials)

	r.PUT("/api/user", auth.ZincAuthMiddleware, handlers.CreateUpdateUser)
	r.DELETE("/api/user/:userID", auth.ZincAuthMiddleware, handlers.DeleteUser)
	r.GET("/api/users", auth.ZincAuthMiddleware, handlers.GetUsers)

	r.PUT("/api/index", auth.ZincAuthMiddleware, handlers.CreateIndex)
	r.GET("/api/index", auth.ZincAuthMiddleware, handlers.ListIndexes)

	r.PUT("/api/:target/document", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.POST("/api/:target/_search", auth.ZincAuthMiddleware, handlers.SearchIndex)

	// elastic compatible APIs
	// Document APIs - https://www.elastic.co/guide/en/elasticsearch/reference/current/docs.html
	// Single document APIs

	// Index - https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-index_.html
	r.PUT("/es/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)

	r.DELETE("/es/:target/_doc/:id", auth.ZincAuthMiddleware, handlers.DeleteDocument)

	r.POST("/es/:target/_doc", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.PUT("/es/:target/_create/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)
	r.POST("/es/:target/_create/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)

	// Update
	r.POST("/es/:target/_update/:id", auth.ZincAuthMiddleware, handlers.UpdateDocument)

	// Bulk update/insert

	r.POST("/es/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)
	r.POST("/es/:target/_bulk", auth.ZincAuthMiddleware, handlers.BulkHandler)

	// Stupid Compatibility
	r.GET("/es/", handlers.Ping)
	r.GET("/es/_xpack", handlers.XPackHandler)
	r.GET("/es/_license", handlers.License)
	r.HEAD("/es/_index_template/:name", handlers.IndexTemplate)
	r.GET("/es/_index_template/:name", handlers.IndexTemplate)
	r.GET("/es/_cat/templates/:name", handlers.CatTemplate)
	r.PUT("/es/_ingest/pipeline/:name", handlers.IngestPipeline)

	// _search
	// hey -n 200 -c 20 -m POST -H "Content-Type: application/json" -d '{"author": "kimchy", "text": "Zincsearch: cool. bonsai cool."}' http://localhost:4080/boos/_doc

	// r.POST("/:target/_search", handlers.SearchIndex)

}
