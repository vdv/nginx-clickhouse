package nginx

import (
	"log"
	"testing"

	"github.com/vdv/nginx-clickhouse/nginx"
)

func TestParseField(t *testing.T) {
	res := nginx.ParseField("bytes_sent", "-")
	log.Printf("bytes_sent was - now ", res)

	res = nginx.ParseField("request_time", "-")
	log.Printf("request_time was - now ", res)
}
