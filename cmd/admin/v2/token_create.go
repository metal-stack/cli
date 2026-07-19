package v2

import (
	"fmt"
	"time"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/cli/cmd/config"
	"github.com/metal-stack/metal-lib/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"
)

type adminToken struct {
	c *config.Config
}

func newAdminTokenCreateCmd(c *config.Config) *cobra.Command {
	w := &adminToken{
		c: c,
	}

	createCmd := &cobra.Command{
		Use:   "token-create",
		Short: "create a token for any user (admin only)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return w.create()
		},
	}

	createCmd.Flags().String("user", "", "the user to create the token for")
	createCmd.Flags().String("description", "", "a short description for the intention to use this token for")
	createCmd.Flags().Duration("expires", 8*time.Hour, "the duration how long the api token is valid")

	return createCmd
}

func (c *adminToken) create() error {
	ctx, cancel := c.c.NewRequestContext()
	defer cancel()

	req := &adminv2.TokenServiceCreateRequest{
		User: pointer.PointerOrNil(viper.GetString("user")),
		TokenCreateRequest: &apiv2.TokenServiceCreateRequest{
			Description: viper.GetString("description"),
			Expires:     durationpb.New(viper.GetDuration("expires")),
		},
	}

	resp, err := c.c.Client.Adminv2().Token().Create(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	_, _ = fmt.Fprintf(c.c.Out, "Make sure to copy your personal access token now as you will not be able to see this again.\n\n")
	_, _ = fmt.Fprintln(c.c.Out, resp.GetSecret())
	_, _ = fmt.Fprintln(c.c.Out)

	return c.c.DescribePrinter.Print(resp.Token)
}
