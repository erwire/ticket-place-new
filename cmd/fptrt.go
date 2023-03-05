package main

import "fptr/pkg/fptr10"

func main() {
	safe, err := fptr10.NewSafe()
	if err != nil {

	}

	safe.Open()
	safe.SetParam(fptr10.LIBFPTR_PARAM_SHIFT_STATE, fptr10.LIBFPTR_SS_EXPIRED)
}
