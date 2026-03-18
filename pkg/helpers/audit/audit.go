package helpersaudit

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RelativeDateTime(s string) (*timestamppb.Timestamp, error) {
	if s == "" {
		return nil, nil
	}
	duration, err := time.ParseDuration(s)
	if err == nil {
		return timestamppb.New(time.Now().Add(-duration)), nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return timestamppb.Now(), fmt.Errorf("failed to convert time: %w", err)
	}
	return timestamppb.New(t), nil
}

func ToPhase(phase string) *apiv2.AuditPhase {
	p, ok := apiv2.AuditPhase_value[phase]
	if !ok {
		return nil
	}

	return new(apiv2.AuditPhase(p))
}

func TryPrettifyBody(trace *apiv2.AuditTrace) *apiv2.AuditTrace {
	if trace.Body != nil {
		trimmed := strings.Trim(*trace.Body, `"`)
		body := map[string]any{}
		if err := json.Unmarshal([]byte(trimmed), &body); err == nil {
			if pretty, err := json.MarshalIndent(body, "", "    "); err == nil {
				trace.Body = new(string(pretty))
			}
		}
	}

	return trace
}
