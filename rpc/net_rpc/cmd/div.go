// Service.Method
package cmd

import "errors"

type DemoService struct {
}

type DivArgs struct {
	A, B int
}

func (d *DemoService) Div(args DivArgs, result *float64) error {
	if args.B == 0 {
		return errors.New("divide has zero value")
	}
	*result = float64(args.A) / float64(args.B)
	return nil
}
