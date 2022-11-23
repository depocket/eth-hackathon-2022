package route

import (
	"depocket.io/app/handler"
	"depocket.io/app/repo"
	"depocket.io/app/service"
	"github.com/dgraph-io/dgo/v200"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NewRoute creates a new router services
func NewRoute(db *gorm.DB, Log *zap.Logger, R *gin.Engine, dgraph *dgo.Dgraph) *Route {
	return &Route{db, Log, R, dgraph}
}

// Route lets us bind specific services when setting up routes
type Route struct {
	DB     *gorm.DB
	Log    *zap.Logger
	R      *gin.Engine
	Dgraph *dgo.Dgraph
}

// Setup instances various repos and services and sets up the routers
func (s *Route) Setup() {
	v1Router := s.R.Group("/v1")

	// repo layer
	dgraphRepo := repo.NewRepoDgraph(s.Dgraph)

	// service layer
	syncAddrService := service.NewSyncAddressService(s.Log, dgraphRepo)
	addrService := service.NewAddressService(s.Log, dgraphRepo)

	// handler
	handler.NewMigrationHandler(v1Router, s.Log, s.Dgraph)
	handler.NewSyncAddressHandler(v1Router, s.Log, s.DB, syncAddrService)
	handler.NewAddressHandler(v1Router, s.Log, s.DB, addrService)
}
