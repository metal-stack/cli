package helpers

import (
	"os"

	metalssh "github.com/metal-stack/metal-lib/pkg/ssh"
)

// sshClient opens an interactive ssh session to the host on port with user, authenticated by the key.
func SShClient(user, keyfile, host string, port int, idToken, project string) error {
	var opts []metalssh.ConnectOpt

	opts = append(opts, metalssh.ConnectOptOutputPassword(idToken))

	if keyfile == "" {
		var err error
		keyfile, err = SearchSSHKey()
		if err != nil {
			return err
		}
	}

	privateKey, err := os.ReadFile(keyfile)
	if err != nil {
		return err
	}

	opts = append(opts, metalssh.ConnectOptOutputPrivateKey(privateKey))

	s, err := metalssh.NewClient(user, host, port, opts...)
	if err != nil {
		return err
	}

	env := map[string]string{
		"LC_METAL_STACK_OIDC_TOKEN": idToken,
		"LC_METAL_STACK_PROJECT":    project,
	}

	sshEnv := metalssh.Env(env)

	return s.Connect(&sshEnv)
}
