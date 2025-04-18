package inits

import (
	"debank_checker_v3/global"
	"debank_checker_v3/util"
)

func InitClients() {
	if len(util.Proxies) == 0 {
		global.Clients = append(global.Clients, util.CreateClient(""))
	} else {
		for _, proxy := range util.Proxies {
			global.Clients = append(global.Clients, util.CreateClient(proxy))
		}
	}
}
