package app

import (
	"depocket.io/app/route"
	"fmt"
	"github.com/apex/gateway"
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
)

type Server struct{}

func inLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func (server *Server) Run(env string) error {
	r := gin.Default()
	log, _ := zap.NewDevelopment()
	dgraphAddress := os.Getenv("DGRAPH_ADDRESS")
	dgraphPort := os.Getenv("DGRAPH_PORT")
	clientConn, err := grpc.Dial(dgraphAddress+":"+dgraphPort, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(clientConn))

	s := route.NewRoute(log, r, dgraphClient)
	s.Setup()

	if inLambda() {
		fmt.Println("running aws lambda in aws")
		return gateway.ListenAndServe(":8080", r)
	} else {
		fmt.Println("running aws lambda in local")
		return r.Run(":" + os.Getenv("SERVER_PORT"))
	}
}
