package route

import (
	"depocket.io/app/handler"
	"depocket.io/app/repo"
	"depocket.io/app/service"
	"github.com/dgraph-io/dgo/v200"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

// NewRoute creates a new router services
func NewRoute(Log *zap.Logger, R *gin.Engine, dgraph *dgo.Dgraph) *Route {
	return &Route{Log, R, dgraph}
}

// Route lets us bind specific services when setting up routes
type Route struct {
	Log    *zap.Logger
	R      *gin.Engine
	Dgraph *dgo.Dgraph
}

// Setup instances various repos and services and sets up the routers
func (s *Route) Setup() {
	v1Router := s.R.Group("/v1")
	v1Router.GET("/hello", func(context *gin.Context) {
		context.JSON(200, gin.H{"status": "ok"})
	})

	v1Router.GET("/ping", func(context *gin.Context) {
		res, err := http.Get("https://google.com")
		if err != nil {
			context.JSON(400, gin.H{"status": "error", "body": err})
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		context.JSON(200, gin.H{"response": bodyString})
	})

	// repo layer
	dgraphRepo := repo.NewRepoDgraph(s.Dgraph)

	// service layer
	syncAddrService := service.NewSyncAddressService(s.Log, dgraphRepo)
	addrService := service.NewAddressService(s.Log, dgraphRepo)

	// handler
	handler.NewMigrationHandler(v1Router, s.Log, s.Dgraph)
	handler.NewSyncAddressHandler(v1Router, s.Log, syncAddrService)
	handler.NewAddressHandler(v1Router, s.Log, addrService)
}
