package ports

type UserRMQ interface {
	Publish(interface{}) error
}
