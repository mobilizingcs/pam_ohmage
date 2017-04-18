package main

import (
  "errors"
)

func isUserAuthenticated( ohmage_url string, username string, password string ) ( bool, error ) {
  // todo: validate username and password strings
  if ohmage_url == "" || username == "" || password == "" {
    return false, errors.New( "Invalid/empty login parameters" )
  }

  return true, nil
}