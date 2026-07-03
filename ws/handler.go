package ws

type handler interface {
	onMessage(message frame)
	onError()
	onPing()
}
