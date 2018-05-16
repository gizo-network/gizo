package p2p

import (
	"net/url"
)

func ParseAddr(addr string) (map[string]interface{}, error) {
	temp := make(map[string]interface{})
	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	temp["pub"] = parsed.User.Username()
	temp["ip"] = parsed.Hostname()
	temp["port"] = parsed.Port()
	return temp, nil
}
