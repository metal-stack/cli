package testresources

import apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

var (
	Version = func() *apiv2.Version {
		return &apiv2.Version{
			Version:   "v0.1.8",
			Revision:  "tags/v0.1.8-0-g476edc0",
			GitSha1:   "477edc0b",
			BuildDate: "2026-03-21T15:35:07+00:00",
		}
	}
)
