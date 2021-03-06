package influxdbw

import (
	"net/http"
	"net/url"

	"github.com/dyweb/gommon/errors"

	"github.com/libtsdb/libtsdb-go/libtsdb/client/genericw"
	"github.com/libtsdb/libtsdb-go/libtsdb/common/influxdb"
	"github.com/libtsdb/libtsdb-go/libtsdb/config"
)

func New(cfg config.InfluxdbClientConfig) (*genericw.HttpClient, error) {
	u, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse server address")
	}
	baseReq, err := http.NewRequest("POST", u.String()+"/write", nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't create base query")
	}
	params := baseReq.URL.Query()
	params.Set("db", cfg.Database)
	baseReq.URL.RawQuery = params.Encode()
	baseReq.Header.Set("User-Agent", "libtsdb")
	c := genericw.NewHttp(influxdb.Meta(), influxdb.NewEncoder(), baseReq)
	return c, nil
}
