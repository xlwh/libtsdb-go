package libtsdb

import (
	"net/http"

	pb "github.com/libtsdb/libtsdb-go/libtsdb/libtsdbpb"
)

// TSBClient returns meta of the database including its protocol, data type support
type TSDBClient interface {
	Meta() Meta
}

type WriteClient interface {
	WriteIntPoint(*pb.PointIntTagged)
	WriteDoublePoint(*pb.PointDoubleTagged)
	WriteSeriesIntTagged(p *pb.SeriesIntTagged)
	WriteSeriesDoubleTagged(p *pb.SeriesDoubleTagged)
	Flush() error
}

type TracedHttpClient interface {
	EnableHttpTrace()
	DisableHttpTrace()
	Trace() HttpTrace
}

type HttpClient interface {
	SetHttpClient(client *http.Client)
	HttpStatusCode() int
}

type HttpWriteClient interface {
	WriteClient
	HttpClient
}

// TODO: figure out the interface for read request...
type ReadClient interface {
	CreateDatabase(db string) error
}
