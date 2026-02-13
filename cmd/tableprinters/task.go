package tableprinters

import (
	"time"

	"github.com/google/uuid"
	"github.com/metal-stack/api/go/enum"
	adminv2 "github.com/metal-stack/api/go/metalstack/admin/v2"
)

func (t *TablePrinter) TaskTable(data []*adminv2.TaskInfo, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "Queue", "When", "Type", "State"}

	if wide {
		header = []string{"ID", "Queue", "When", "Type", "State", "Issued At", "Payload", "Result"}
	}
	for _, task := range data {
		id := task.Id
		queue := task.Queue
		typeString := task.Type
		state, err := enum.GetStringValue(task.State)
		if err != nil {
			state = new("unknown")
		}
		payload := string(task.Payload)
		result := string(task.Result)

		parsed, err := uuid.Parse(id)
		if err != nil {
			return nil, nil, err
		}

		sec, nano := parsed.Time().UnixTime()
		issuedAt := time.Unix(sec, nano)
		when := humanizeDuration(time.Since(issuedAt))

		if wide {
			rows = append(rows, []string{id, queue, when, typeString, *state, issuedAt.String(), payload, result})
		} else {
			rows = append(rows, []string{id, queue, when, typeString, *state})
		}
	}

	t.t.DisableAutoWrap(false)

	return header, rows, nil
}
