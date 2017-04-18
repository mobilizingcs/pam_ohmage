package main

/*
#cgo LDFLAGS: -lpam -fPIC
#include <security/pam_appl.h>
#include <stdlib.h>
#include <string.h>
char *string_from_argv(int i, char **argv);
char *get_username(pam_handle_t *pamh);
char *get_password(pam_handle_t *pamh);
*/
import "C"

import (
  "fmt"
  "strings"
  "unsafe"
  "errors"
  "net/url"
)

//export pam_sm_authenticate
func pam_sm_authenticate( pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char ) C.int {
  _username := C.get_username( pamh )
  if _username == nil {
    return C.PAM_USER_UNKNOWN
  }
  defer C.free( unsafe.Pointer( _username ) )
  username := C.GoString(_username);

  _password := C.get_password( pamh )
  if _password == nil {
    return C.PAM_ABORT
  }
  defer C.free( unsafe.Pointer( _password ) )
  password := C.GoString(_password);

  cli_params := sliceFromArgv( argc, argv )
  ohmage_url, err := parseOhmageUrl( cli_params[0] )
  if err != nil {
    fmt.Println( err )
    return C.PAM_ABORT
  }

  authenticated, err := isUserAuthenticated( ohmage_url, username, password )
  if err != nil {
    fmt.Println( err )
    return C.PAM_ABORT
  } else if authenticated {
    return C.PAM_SUCCESS
  } else {
    return C.PAM_USER_UNKNOWN
  }
}

func parseOhmageUrl( _parameter string ) ( string, error ) {
  parameter := strings.Split( _parameter, "=" )
  if parameter[ 0 ] == "url" {
    parsed_url, err := url.ParseRequestURI( parameter[ 1 ] )
    if( err != nil  ) {
      return "", errors.New( "invalid ohmage url" )
    } else {
      return parsed_url.String( ), nil
    }
  } else {
    return "", errors.New( "url parameter not found" )
  }
}

func sliceFromArgv( argc C.int, argv **C.char ) [ ]string {
  result := make( []string, 0, argc )
  for i := 0; i < int( argc ); i++ {
    str := C.string_from_argv( C.int( i ), argv )
    defer C.free( unsafe.Pointer( str ) )
    result = append( result, C.GoString( str ) )
  }
  return result
}

func main( ) { }
