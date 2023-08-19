package sender

type Sender interface {
	Send(filename, title, description string) (string, error)
}
