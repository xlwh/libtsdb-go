package influxdbw

import (
	"testing"

	asst "github.com/stretchr/testify/assert"

	"github.com/libtsdb/libtsdb-go/libtsdb/config"
	pb "github.com/libtsdb/libtsdb-go/libtsdb/libtsdbpb"
)

// TODO: add flag to toggle test base on environ variable ... maybe testutil to gommon, travis etc.
func TestClient_WriteIntPoint(t *testing.T) {
	t.Skip("requires influxdb running")

	assert := asst.New(t)
	c, err := New(*config.NewInfluxdbClientConfig())
	c.EnableHttpTrace()
	assert.Nil(err)
	// TODO: util for point generator
	p := &pb.PointIntTagged{
		Name:  "temperature",
		Point: pb.PointInt{T: int64(1434055562000000035), V: 35},
		Tags: []pb.Tag{
			{K: "machine", V: "unit42"},
			{K: "type", V: "assembly"},
		},
	}
	c.WriteIntPoint(p)
	err = c.Flush()
	assert.Nil(err)
	trace := c.Trace()
	assert.Equal(204, trace.GetCode())
	msize := len(p.Name)
	for _, tg := range p.Tags {
		msize += len(tg.K) + len(tg.V)
	}
	t.Logf("payload size %d", trace.GetPayloadSize())
	t.Logf("raw size %d", trace.GetRawSize())
	t.Log("raw meta size", trace.GetRawMetaSize())
	assert.Equal(msize, trace.GetRawMetaSize())
	//t.Logf("%v", trace)
}

func TestClient_WriteDoublePoint(t *testing.T) {
	t.Skip("requires influxdb running")

	assert := asst.New(t)
	c, err := New(*config.NewInfluxdbClientConfig())
	assert.Nil(err)
	// TODO: influxdb even allow different type in a same series?
	c.WriteDoublePoint(&pb.PointDoubleTagged{
		Name:  "temperatured",
		Point: pb.PointDouble{T: int64(1434055562000000036), V: 35.132},
		Tags: []pb.Tag{
			{K: "machine", V: "unit42"},
			{K: "type", V: "assembly"},
		},
	})
	err = c.Flush()
	assert.Nil(err)
}
