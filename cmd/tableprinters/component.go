package tableprinters

import (
	"time"

	"github.com/fatih/color"
	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) ComponentTable(data []*apiv2.Component, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Type", "Identifier", "Started", "Age", "Version", "Token", "Token Expires In"}
	)

	for _, c := range data {
		typeString, err := enum.GetStringValue(c.Type)
		if err != nil {
			return nil, nil, err
		}

		started := humanizeDuration(time.Since(c.StartedAt.AsTime()))
		ageAsDuration := time.Since(c.ReportedAt.AsTime())
		age := humanizeDuration(ageAsDuration)

		if c.Interval != nil && c.Interval.AsDuration() < ageAsDuration {
			age = color.RedString(age)
		}

		tokenExpirationAsDuration := time.Until(c.Token.Expires.AsTime())
		tokenExpiresIn := humanizeDuration(tokenExpirationAsDuration)

		if tokenExpirationAsDuration < time.Hour {
			tokenExpiresIn = color.YellowString(tokenExpiresIn)
		} else if tokenExpirationAsDuration < 0 {
			tokenExpiresIn = color.RedString(tokenExpiresIn)
		}

		rows = append(rows, []string{c.Uuid, *typeString, c.Identifier, started, age, c.Version.Version, c.Token.Uuid, tokenExpiresIn})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
