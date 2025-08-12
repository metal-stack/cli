package tableprinters

import (
	"strings"
	"time"

	"github.com/metal-stack/api/go/enum"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"github.com/metal-stack/metal-lib/pkg/pointer"
)

func (t *TablePrinter) ImageTable(data []*apiv2.Image, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"ID", "Name", "Description", "Features", "Expiration", "Status"}
	)

	for _, image := range data {
		var (
			features []string
		)

		for _, f := range image.Features {
			feature, err := enum.GetStringValue(f)
			if err != nil {
				return nil, nil, err
			}
			features = append(features, *feature)
		}

		classification, err := enum.GetStringValue(image.Classification)
		if err != nil {
			return nil, nil, err
		}

		var expiresIn string
		if image.ExpiresAt != nil {
			expiresIn = humanizeDuration(time.Until(image.ExpiresAt.AsTime()))
		}

		rows = append(rows, []string{image.Id, pointer.SafeDeref(image.Name), pointer.SafeDeref(image.Description), strings.Join(features, ","), expiresIn, *classification})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
