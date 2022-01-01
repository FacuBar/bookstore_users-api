package ports

type UserRMQ interface {
	Publish(string, interface{}) error
}
