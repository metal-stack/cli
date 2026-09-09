package helpers

import (
	"os"

	consoleapi "github.com/metal-stack/metal-console/api"
	metalssh "github.com/metal-stack/metal-lib/pkg/ssh"
)

// SShClient opens an interactive ssh session to the host on port with user, authenticated by the key.
func SShClient(user, keyfile, host string, port int, idToken string, project *string) error {
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
		consoleapi.OidcTokenEnv: idToken,
	}
	if project != nil {
		env[consoleapi.ProjectEnv] = *project
	}

	sshEnv := metalssh.Env(env)

	return s.Connect(&sshEnv)
}
