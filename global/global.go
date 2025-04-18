package global

import (
	"debank_checker_v3/customTypes"

	tls_client "github.com/bogdanfinn/tls-client"
)

var (
	ConfigFile      customTypes.Config
	Clients         []tls_client.HttpClient
	TargetProgress  int
	CurrentProgress int32
)
