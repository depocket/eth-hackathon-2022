package app

import (
	"depocket.io/app/route"
	"depocket.io/pkgs/database/connection"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
)

type Server struct{}

func (server *Server) Run(env string) error {
	r := gin.Default()
	log, _ := zap.NewDevelopment()
	db := connection.NewConnection(zap.NewNop(), nil)
	dgraphAddress := os.Getenv("DGRAPH_ADDRESS")
	dgraphPort := os.Getenv("DGRAPH_PORT")
	clientConn, err := grpc.Dial(dgraphAddress+":"+dgraphPort, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(clientConn))

	s := route.NewRoute(db, log, r, dgraphClient)
	s.Setup()

	return r.Run(":" + os.Getenv("SERVER_PORT"))
}
