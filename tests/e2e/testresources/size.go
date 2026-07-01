package testresources

import apiv2 "github.com/metal-stack/api/go/metalstack/api/v2"

var (
	Size1 = func() *apiv2.Size {
		return &apiv2.Size{
			Id:          "v1-medium-x86",
			Name:        new("v1-medium-x86"),
			Description: new("Virtual size for mini-lab"),
			Meta: &apiv2.Meta{
				Labels: &apiv2.Labels{
					Labels: map[string]string{
						"a": "b",
					},
				},
			},
			Constraints: []*apiv2.SizeConstraint{
				{
					Type: apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_CORES,
					Min:  4,
					Max:  4,
				},
				{
					Type: apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_MEMORY,
					Min:  500000000,
					Max:  4000000000,
				},
				{
					Type: apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_STORAGE,
					Min:  1000000000,
					Max:  100000000000,
				},
			},
		}
	}
	Size2 = func() *apiv2.Size {
		return &apiv2.Size{
			Id:          "g1-medium-x86",
			Name:        new("g1-medium-x86"),
			Description: new("A medium sized GPU server"),
			Meta:        &apiv2.Meta{},
			Constraints: []*apiv2.SizeConstraint{
				{
					Identifier: new("Intel(R) Xeon(R) Gold 6426Y"),
					Type:       apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_CORES,
					Min:        32,
					Max:        32,
				},
				{

					Type: apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_MEMORY,
					Min:  274877906944,
					Max:  274877906944,
				},
				{
					Identifier: new("/dev/*"),
					Type:       apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_STORAGE,
					Min:        1840378724352,
					Max:        1840378724352,
				},
				{
					Identifier: new("AD102GL [RTX 6000 Ada Generation]"),
					Type:       apiv2.SizeConstraintType_SIZE_CONSTRAINT_TYPE_GPU,
					Min:        1,
					Max:        1,
				},
			},
		}
	}
)
