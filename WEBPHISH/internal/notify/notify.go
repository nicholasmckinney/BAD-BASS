package notify

type Notifier interface {
	Notify() bool
}
