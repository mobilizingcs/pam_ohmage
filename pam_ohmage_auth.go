package main

import (
  "github.com/pkg/errors"
  "github.com/mobilizingcs/gohmage"
)

func isUserAuthenticated( ohmage_url string, username string, password string ) ( bool, error ) {
  log.Debug( "authenticate module called" )
  if ohmage_url == "" || username == "" || password == "" {
    log.Debug( "Invalid/empty login parameters" )
    return false, errors.New( "Invalid/empty login parameters" )
  }
  ohmage := gohmage.NewClient( ohmage_url, "pam_ohmage" )
  auth_status, err := ohmage.UserAuthToken( username, password )
  if err != nil || auth_status == false {
    log.Debug( "Ohmage authentication call failed" )
    return false, err
  }
  user_info, err := ohmage.UserInfoRead( )
  if err != nil {
    log.Debug( "Ohmage user info call failed" )
    return false, errors.Wrap( err, "Unable to fetch user information" )
  }
  user_classes, err := user_info.GetObject( "data", username, "classes" )
  if err != nil {
    log.Debug( "Failed to find user classes" )
    return false, errors.Wrap( err, "Unable to find user classes" )
  }
  // todo: pick up test class from PAM configuration
  test_class, err := user_classes.GetString( "urn:class:rstudio" )
  if err != nil || test_class == "" {
    log.Debug( "Failed to find test class among user's classes" )
    return false, errors.Wrap( err, "Test class not among the classes user is participant of" )
  }
  return true, nil
}