package main

import (
	"fmt"
	"fptr/pkg/fptr10"
	"log"
)

func main() {
	fptr, err := fptr10.NewSafe()
	if err != nil {
		log.Fatal(err.Error())
	}

	fptr.Open()
	//CloseShift(fptr)

	fptr.SetParam(1021, "Егор Ишуткин")
	fptr.SetParam(1203, "500100732259")
	fptr.OperatorLogin()

	defer fptr.Destroy()
	fmt.Println(fptr.IsOpened())
	fmt.Println(fptr.GetParamInt(fptr10.LIBFPTR_PARAM_SHIFT_STATE) == fptr10.LIBFPTR_SS_CLOSED)

}

func CloseShift(fptr *fptr10.IFptr) {
	fptr.SetParam(fptr10.LIBFPTR_PARAM_REPORT_TYPE, fptr10.LIBFPTR_RT_CLOSE_SHIFT)
	fptr.Report()
}
