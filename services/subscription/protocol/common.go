package protocol

import (
	"fmt"

	"github.com/xtls/xray-core/infra/conf"
)

// TODO: IDN encoding
func getHost(address *conf.Address, port uint16) (string, error) {
	if address.String() == "" {
		return "", ErrNoHost
	}
	if !(port >= 1 && port <= 65535) {
		return "", IllegalPortError(port)
	}
	return fmt.Sprintf("%v:%v", address, port), nil
}
