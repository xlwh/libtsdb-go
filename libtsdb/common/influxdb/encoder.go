package influxdb

import (
	"strconv"

	"github.com/libtsdb/libtsdb-go/libtsdb/common"
	pb "github.com/libtsdb/libtsdb-go/libtsdb/libtsdbpb"
	"github.com/libtsdb/libtsdb-go/libtsdb/util/bytesutil"
)

const (
	defaultField = "v"
)

var _ common.Encoder = (*Encoder)(nil)

// Encoders encodes points with tags into InfluxDB's line protocol, it ONLY has a single field called v
//
// ref https://github.com/influxdata/influxdb/blob/master/models/points.go#L2267 appendField
type Encoder struct {
	bytesutil.Buffer
	// DefaultField is used when encoding single field series, which is the case for most TSDB but not InfluxDB
	// other tsdb (name, tags, value, ts)     : cpu.usage host=i7szx,dc=us-east 1.0 1359788400000
	// influxdb (name, tags, field=value, ts) : cpu.usage host=i7szx,dc=us-east v=1.0 1359788400000
	DefaultField string
}

func NewEncoder() *Encoder {
	return &Encoder{
		DefaultField: defaultField,
	}
}

// temperature,machine=unit42,type=assembly internal=32,external=100 1434055562000000035
func (e *Encoder) WritePointIntTagged(p *pb.PointIntTagged) {
	e.Buf = append(e.Buf, p.Name...)
	e.Buf = append(e.Buf, ',')
	for _, tag := range p.Tags {
		e.Buf = append(e.Buf, tag.K...)
		e.Buf = append(e.Buf, '=')
		e.Buf = append(e.Buf, tag.V...)
		e.Buf = append(e.Buf, ',')
	}
	e.Buf[len(e.Buf)-1] = ' '
	e.Buf = append(e.Buf, e.DefaultField...)
	e.Buf = append(e.Buf, '=')
	e.Buf = strconv.AppendInt(e.Buf, p.Point.V, 10)
	e.Buf = append(e.Buf, ' ')
	e.Buf = strconv.AppendInt(e.Buf, p.Point.T, 10)
	e.Buf = append(e.Buf, '\n')
}

func (e *Encoder) WritePointDoubleTagged(p *pb.PointDoubleTagged) {
	e.Buf = append(e.Buf, p.Name...)
	e.Buf = append(e.Buf, ',')
	for _, tag := range p.Tags {
		e.Buf = append(e.Buf, tag.K...)
		e.Buf = append(e.Buf, '=')
		e.Buf = append(e.Buf, tag.V...)
		e.Buf = append(e.Buf, ',')
	}
	e.Buf[len(e.Buf)-1] = ' '
	e.Buf = append(e.Buf, e.DefaultField...)
	e.Buf = append(e.Buf, '=')
	// TODO: most part are copy and pasted except this line ...
	e.Buf = strconv.AppendFloat(e.Buf, p.Point.V, 'f', -1, 64)
	e.Buf = append(e.Buf, ' ')
	e.Buf = strconv.AppendInt(e.Buf, p.Point.T, 10)
	e.Buf = append(e.Buf, '\n')
}

func (e *Encoder) WriteSeriesIntTagged(p *pb.SeriesIntTagged) {
	// NOTE: InfluxDB does not support ingest multiple points in one line
	// first write the key, measurement + tags + default field
	var header []byte
	header = append(header, p.Name...)
	header = append(header, ',')
	for _, tag := range p.Tags {
		header = append(header, tag.K...)
		header = append(header, '=')
		header = append(header, tag.V...)
		header = append(header, ',')
	}
	header[len(header)-1] = ' '
	header = append(header, e.DefaultField...)
	header = append(header, '=')
	// then duplicate the header in every line
	for i := range p.Points {
		e.Buf = append(e.Buf, header...)
		e.Buf = strconv.AppendInt(e.Buf, p.Points[i].V, 10)
		e.Buf = append(e.Buf, ' ')
		e.Buf = strconv.AppendInt(e.Buf, p.Points[i].T, 10)
		e.Buf = append(e.Buf, '\n')
	}
}

func (e *Encoder) WriteSeriesDoubleTagged(p *pb.SeriesDoubleTagged) {
	// NOTE: InfluxDB does not support ingest multiple points in one line
	// first write the key, measurement + tags + default field
	var header []byte
	header = append(header, p.Name...)
	header = append(header, ',')
	for _, tag := range p.Tags {
		header = append(header, tag.K...)
		header = append(header, '=')
		header = append(header, tag.V...)
		header = append(header, ',')
	}
	header[len(header)-1] = ' '
	header = append(header, e.DefaultField...)
	header = append(header, '=')
	// then duplicate the header in every line
	for i := range p.Points {
		e.Buf = append(e.Buf, header...)
		e.Buf = strconv.AppendFloat(e.Buf, p.Points[i].V, 'f', -1, 64)
		e.Buf = append(e.Buf, ' ')
		e.Buf = strconv.AppendInt(e.Buf, p.Points[i].T, 10)
		e.Buf = append(e.Buf, '\n')
	}
}
