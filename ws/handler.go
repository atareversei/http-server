package ws

type handler interface {
	onMessage(message Frame)
	onError()
	onPing()
}
