package bus

type EventHandler interface {
	Apply(event Event)
}
