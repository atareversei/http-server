package http

const (
	HeaderKeyUpgrade    = "upgrade"
	HeaderKeyConnection = "connection"
	HeaderKeyOrigin = "origin"

	HeaderValueConnection_Upgrade = "Upgrade"

	// WebSockets
	HeaderKeyWSVersion   = "sec-webSocket-version"
	HeaderValueWSUpgrade = "websocket"
	HeaderValueWSMethod  = MethodGet
	HeaderValueWSVersion = "13"
)
