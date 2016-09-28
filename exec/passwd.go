package exec

import "github.com/markuskobler/authorized/provider"

func CreateUser(u provider.User) (err error) {

	return
}

func AuthorizeSSHKeys(u provider.User) (err error) {
	if len(u.SSHAuthorizedKeys) == 0 {
		return nil
	}

	return
}
