package main

// #cgo LDFLAGS: -lpam -fPIC
// #include <security/pam_appl.h>
// #include <stdlib.h>
import "C"

import (
  "fmt"
)

//export pam_sm_authenticate
func pam_sm_authenticate( pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char ) C.int {
  fmt.Printf( "Called\n" );
  return C.PAM_SUCCESS
}

func main( ) { }
