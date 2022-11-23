package repo

import (
	"github.com/dgraph-io/dgo/v200"
	"github.com/dgraph-io/dgo/v200/protos/api"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestDgraph_FullFlow(t *testing.T) {
	godotenv.Load("../../.env")
	dgraphAddress := os.Getenv("DGRAPH_ADDRESS")
	dgraphPort := os.Getenv("DGRAPH_PORT")
	clientConn, _ := grpc.Dial(dgraphAddress+":"+dgraphPort, grpc.WithInsecure())
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(clientConn))

	from, _ := time.Parse(time.RFC3339, "2020-09-01T10:16:28+07:00")
	to, _ := time.Parse(time.RFC3339, "2020-09-20T10:16:28+07:00")
	type fields struct {
		DB    *dgo.Dgraph
		debug bool
	}
	type args struct {
		depth   int
		address string
		token   string
		from    time.Time
		to      time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "dgraphClient", fields: fields{
			DB:    dgraphClient,
			debug: false,
		}, args: struct {
			depth   int
			address string
			token   string
			from    time.Time
			to      time.Time
		}{
			depth:   10,
			address: "0x9fb34d03374786a14e776d246f62eabdd9caaefe",
			token:   "0xe9e7cea3dedca5984780bafc599bd69add087d56",
			from:    from,
			to:      to},
			want:    nil,
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Dgraph{
				DB:    tt.fields.DB,
				debug: tt.fields.debug,
			}
			got, err := r.FullFlow(tt.args.depth, tt.args.address, tt.args.token, tt.args.from, tt.args.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("FullFlow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FullFlow() got = %v, want %v", got, tt.want)
			}
		})
	}
}
