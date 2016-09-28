package exec

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/markuskobler/authorized/provider"
)

const (
	authorizedKeysFile = "authorized_keys"
	authorizedKeysDir  = "authorized_keys.d"
	preservedKeysName  = "orig_authorized_keys"
	sshDir             = ".ssh"

	lockFile  = ".authorized_keys.d.lock"       // In "~/".
	stageFile = ".authorized_keys.d.stage_file" // In "~/.ssh".
	stageDir  = ".authorized_keys.d.stage_dir"  // In "~/.ssh".
)

// Open opens the authorized keys directory for the supplied user.
// If create is false, Open will fail if no directory exists yet.
// If create is true, Open will create the directory if it doesn't exist,
// preserving the authorized_keys file in the process.
// After a successful open, Close should be called when finished to unlock
// the directory.
func openUser(usr *provider.User, create bool) (*sshAuthorizedKeysDir, error) {
	l, err := acquireLock(usr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			l.Close()
		}
	}()

	akd, err := opendir(authKeysDirPath(usr))
	if err != nil && (!create || !os.IsNotExist(err)) {
		return nil, err
	} else if os.IsNotExist(err) {
		akd, err = createAuthorizedKeysDir(usr)
		if err != nil {
			return nil, err
		}
	}

	akd.lock = l
	akd.user = usr
	return akd, nil
}

// acquireLock locks the lock file for the given user's authorized_keys.d.
// A lock file is created if it doesn't already exist.
// The locking is currently a simple coarse-grained mutex held for the
// Open()-Close() duration, implemented using a lock file in the user's ~/.
func acquireLock(u *provider.User) (*os.File, error) {
	f, err := as_user.OpenFile(u, lockFilePath(u),
		syscall.O_CREAT|syscall.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

// lockFilePath returns the path to the lock file for the user.
func lockFilePath(u *provider.User) string {
	return filepath.Join(u.HomeDir, lockFile)
}

// opendir opens the authorized keys directory.
func opendir(dir string) (*sshAuthorizedKeysDir, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", dir)
	}
	return &sshAuthorizedKeysDir{path: dir}, nil
}

// SSHAuthorizedKeysDir represents an opened user's authorized_keys.d.
type sshAuthorizedKeysDir struct {
	path string         // Path to authorized_keys.d directory.
	user *provider.User // User of the directory.
	lock *os.File       // Lock file for serializing Open()-Close().
}

// authKeysFilePath returns the path to the authorized_keys file for the user.
func authKeysFilePath(u *user.User) string {
	return filepath.Join(sshDirPath(u), authorizedKeysFile)
}

// sshDirPath returns the path to the .ssh dir for the user.
func sshDirPath(u *user.User) string {
	return filepath.Join(u.HomeDir, sshDir)
}
