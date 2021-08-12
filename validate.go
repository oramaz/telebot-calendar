package calendar

import (
	"errors"

	vd "github.com/go-ozzo/ozzo-validation"
)

func (opt *Options) validate() error {
	return vd.ValidateStruct(opt,
		vd.Field(&opt.YearRange, vd.Required, vd.By(func(v interface{}) error {
			rng := v.([2]int)
			if rng[0] < MinYearLimit || rng[1] > MaxYearLimit {
				return errors.New(errYearRangeValue)
			}
			return nil
		})),
		vd.Field(&opt.InitialYear, vd.Required,
			vd.Min(opt.YearRange[0]),
			vd.Max(opt.YearRange[1]),
		),
		vd.Field(&opt.InitialMonth, vd.Required, vd.Min(1), vd.Max(12)),
	)
}
