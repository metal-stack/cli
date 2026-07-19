package tableprinters

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

func (t *TablePrinter) VPNServiceAuthKeyResponseTable(data []*adminv2.VPNServiceAuthKeyResponse, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Address", "AuthKey", "Ephemeral", "Expires At", "Created At"}
	)

	for _, r := range data {
		expiresAt := ""
		if r.ExpiresAt != nil {
			expiresAt = r.ExpiresAt.AsTime().Format(time.RFC3339)
		}
		createdAt := ""
		if r.CreatedAt != nil {
			createdAt = r.CreatedAt.AsTime().Format(time.RFC3339)
		}

		rows = append(rows, []string{r.Address, r.AuthKey, fmt.Sprintf("%v", r.Ephemeral), expiresAt, createdAt})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func (t *TablePrinter) ImageUsageTable(data []*apiv2.ImageUsage, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Image", "Used By"}
	)

	for _, u := range data {
		imageID := ""
		if u.Image != nil {
			imageID = u.Image.Id
		}
		rows = append(rows, []string{imageID, joinOrEmpty(u.UsedBy)})
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}

func formatInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}

func formatInt32(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func formatUint32(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

func formatUint64(v uint64) string {
	return strconv.FormatUint(v, 10)
}

func joinOrEmpty(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, ", ")
}
