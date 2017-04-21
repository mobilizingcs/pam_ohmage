package main

/*
#cgo LDFLAGS: -lpam -fPIC
#define PAM_SM_AUTH
#define PAM_SM_SESSION
#include <security/pam_modules.h>
#include <stdlib.h>
#include <string.h>
char *string_from_argv(int i, char **argv);
char *get_username(pam_handle_t *pamh);
char *get_password(pam_handle_t *pamh);
*/
import "C"

import (
  "unsafe"
  "strings"
  "github.com/op/go-logging"
)

//export pamAuthenticate
func pamAuthenticate( pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char ) C.int {
  username := getPamUsername( pamh )
  if username == "" {
    return C.PAM_USER_UNKNOWN
  }

  _password := C.get_password( pamh )
  if _password == nil {
    return C.PAM_ABORT
  }
  defer C.free( unsafe.Pointer( _password ) )
  password := C.GoString(_password);
  if password == "" {
    // No ohmage user can have an empty password!
    return C.PAM_USER_UNKNOWN
  }

  cli_params := mapFromArgv( argc, argv )
  ohmage_url, err := parseOhmageUrl( cli_params[ "url" ] )
  if err != nil {
    return C.PAM_ABORT
  }
  if cli_params[ "debug" ] == "true" {
    logging.SetLevel( logging.DEBUG, "pam_ohmage" )
  }

  authenticated, err := isUserAuthenticated( ohmage_url, username, password )

  if err != nil {
    log.Error( "userame:", username, err )
    return C.PAM_AUTH_ERR
  } else if authenticated {
    // RStudio expects the local account to be present after authentication succeeds.
    // This is a mechanism to get around this assumption: We call the open session
    // module from within this module before returning.
    // The open session module is still needed for RStudio PAM sessions
    log.Debug( "Calling open_session module" )
    user_account_ready, err := isUserAccountReady( username )
    if err != nil {
      log.Error( err )
      return C.PAM_ABORT
    } else if user_account_ready {
      return C.PAM_SUCCESS
    } else {
      return C.PAM_ABORT
    }
  } else {
    return C.PAM_AUTH_ERR
  }
}

//export pamOpenSession
func pamOpenSession( pamh *C.pam_handle_t, flags C.int, argc C.int, argv **C.char ) C.int {
  cli_params := mapFromArgv( argc, argv )
  if cli_params[ "debug" ] == "true" {
    logging.SetLevel( logging.DEBUG, "pam_ohmage" )
  }
  username := getPamUsername( pamh )
  if username != "" {
    user_account_ready, err := isUserAccountReady( username )
    if err != nil {
      log.Error( "username:", username, err )
      return C.PAM_SESSION_ERR
    } else if user_account_ready {
      return C.PAM_SUCCESS
    }
  }
  return C.PAM_SESSION_ERR
}

func getPamUsername( pamh *C.pam_handle_t ) string {
  _username := C.get_username( pamh )
  if _username == nil {
    return ""
  }
  defer C.free( unsafe.Pointer( _username ) )
  username := C.GoString(_username);
  return username
}

func mapFromArgv( argc C.int, argv **C.char ) map[string]string {
  parameters := make( []string, 0, argc )
  for i := 0; i < int( argc ); i++ {
    str := C.string_from_argv( C.int( i ), argv )
    defer C.free( unsafe.Pointer( str ) )
    parameters = append( parameters, C.GoString( str ) )
  }
  result := make( map[string]string )
  if len( parameters ) > 0 {
    for _ , _parameter := range parameters {
      parameter := strings.Split( _parameter, "=" )
      if len( parameter ) > 1 && parameter[ 0 ] != "" {
        result[ parameter[ 0 ] ] = parameter[ 1 ]
      }
    }
  }
  return result
}

func main( ) { }
