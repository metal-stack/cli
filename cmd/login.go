package cmd

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/fatih/color"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/genericcli"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
)

type login struct {
	c *config.Config
}

func newLoginCmd(c *config.Config) *cobra.Command {
	w := &login{
		c: c,
	}

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.login()
		},
	}

	loginCmd.Flags().String("provider", "oidc", "the provider used to login with")
	loginCmd.Flags().String("context", "", "the context into which the token gets injected, if not specified it uses the current context or creates a context named default in case there is no current context set")
	loginCmd.Flags().String("admin-role", "", "operators can use this flag to issue an admin token with the token retrieved from login and store this into context")

	genericcli.Must(loginCmd.Flags().MarkHidden("admin-role"))
	genericcli.Must(loginCmd.RegisterFlagCompletionFunc("provider", cobra.FixedCompletions([]string{"oidc"}, cobra.ShellCompDirectiveNoFileComp)))
	genericcli.Must(loginCmd.RegisterFlagCompletionFunc("admin-role", c.Completion.TokenAdminRoleCompletion))

	return loginCmd
}

func (l *login) login() error {
	provider := l.c.GetProvider()
	if provider == "" {
		return errors.New("provider must be specified")
	}

	tokenChan := make(chan string)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		tokenChan <- r.URL.Query().Get("token")

		http.Redirect(w, r, "https://metal-stack.io", http.StatusSeeOther)
	})

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	server := http.Server{Addr: listener.Addr().String(), ReadTimeout: 2 * time.Second}

	go func() {
		fmt.Printf("Starting server at http://%s...\n", listener.Addr().String())
		err = server.Serve(listener) //nolint
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("http server closed unexpectedly: %w", err))
		}
	}()

	url := fmt.Sprintf("%s/auth/%s?redirect-url=http://%s/callback", l.c.GetApiURL(), provider, listener.Addr().String()) // TODO(vknabel): nicify please

	err = exec.Command("xdg-open", url).Run() //nolint // TODO probably broken on MAC?
	if err != nil {
		return fmt.Errorf("error opening browser: %w", err)
	}

	token := <-tokenChan

	err = server.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("unable to close http server: %w", err)
	}
	_ = listener.Close()

	if token == "" {
		return errors.New("no token was retrieved")
	}

	if viper.IsSet("admin-role") {
		mc, err := newApiClient(l.c.GetApiURL(), token)
		if err != nil {
			return err
		}

		tokenResp, err := mc.Apiv2().Token().Create(context.Background(), &apiv2.TokenServiceCreateRequest{
			Description: "admin access issues by metal cli",
			Expires:     durationpb.New(3 * time.Hour),
			AdminRole:   pointer.Pointer(apiv2.AdminRole((apiv2.AdminRole_value[viper.GetString("admin-role")]))),
		})
		if err != nil {
			return fmt.Errorf("unable to issue admin token: %w", err)
		}

		token = tokenResp.Secret
	}

	var ctx *genericcli.Context
	var defaultCtx bool
	name := viper.GetString("context")

	if viper.IsSet("context") {
		ctx, err = l.c.ContextManager.Get(name)
		if err != nil {
			return err
		}
	} else {
		ctx, err = l.c.ContextManager.GetCurrentContext()
		defaultCtx = err != nil || ctx == nil
	}

	fmt.Println(defaultCtx)
	fmt.Println(ctx)

	var project string
	if defaultCtx || ctx.DefaultProject == "" {
		mc, err := newApiClient(l.c.GetApiURL(), token)
		if err != nil {
			return err
		}

		projects, err := mc.Apiv2().Project().List(context.Background(), &apiv2.ProjectServiceListRequest{})
		if err != nil {
			return fmt.Errorf("unable to retrieve project list: %w", err)
		}

		if len(projects.Projects) > 0 {
			project = projects.Projects[0].Uuid
		}
	}

	if defaultCtx {
		ctx, err = l.c.ContextManager.Create(&genericcli.Context{
			Name:           cmp.Or(name, string(genericcli.DefaultContextName)),
			APIURL:         pointer.PointerOrNil((l.c.GetApiURL())),
			APIToken:       token,
			DefaultProject: project,
			// Timeout:        &0,
			Provider:  provider,
			IsCurrent: true,
		})
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintf(l.c.Out, "%s Context \"%s\" is actived \n", color.GreenString("✔"), color.GreenString(ctx.Name))
		_, _ = fmt.Fprintf(l.c.Out, "%s Login successful!\n", color.GreenString("✔"))
		return nil
	}

	_, err = l.c.ContextManager.Update(&genericcli.ContextUpdateRequest{
		Name:           ctx.Name,
		APIURL:         ctx.APIURL,
		APIToken:       &token,
		DefaultProject: pointer.Pointer(cmp.Or(ctx.DefaultProject, project)),
		Provider:       &provider,
		Activate:       true,
	})
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(l.c.Out, "%s Login successful!\n", color.GreenString("✔"))

	return nil
}
