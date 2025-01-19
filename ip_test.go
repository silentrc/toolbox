package toolbox

import (
	"testing"
)

func TestGetIPv6Addr(t *testing.T) {
	NewUtils().NewIP().Init("./conf/GeoLite2-City.mmdb", "./conf/ipv6wry.db")

}
