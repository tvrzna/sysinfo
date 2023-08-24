package main

// #include <stdlib.h>
import "C"

type Loadavg struct {
	Loadavg1  float32
	Loadavg5  float32
	Loadavg15 float32
}

func GetLoadavg() *Loadavg {
	var arr [3]C.double
	result := C.getloadavg(&arr[0], 3)
	if result != 3 {
		return nil
	}
	return &Loadavg{Loadavg1: float32(arr[0]), Loadavg5: float32(arr[1]), Loadavg15: float32(arr[2])}
}
