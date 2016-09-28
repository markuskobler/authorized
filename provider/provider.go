package provider

type Provider interface {
	Users() <-chan []User
}

type User struct {
	Name     string
	SSHKeys  []string
	HomeDir  string
	Shell    string
	Disabled bool
}
