package testresources

import (
	apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"
)

var (
	SizeReservation1 = func() *apiv2.SizeReservation {
		return &apiv2.SizeReservation{
			Id:          "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			Name:        "reservation-a",
			Description: "first reservation",
			Project:     Project1().Uuid,
			Size:        Size1().Id,
			Partitions:  []string{Partition1().Id},
			Amount:      3,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
		}
	}
	SizeReservation2 = func() *apiv2.SizeReservation {
		return &apiv2.SizeReservation{
			Id:          "b2c3d4e5-f6a7-8901-bcde-f12345678901",
			Name:        "reservation-b",
			Description: "second reservation",
			Project:     Project2().Uuid,
			Size:        Size2().Id,
			Partitions:  []string{Partition2().Id},
			Amount:      1,
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"c": "d",
					},
				},
			},
		}
	}
)
