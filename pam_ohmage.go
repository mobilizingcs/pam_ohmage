package main

import (
  "github.com/pkg/errors"
  "gohmage"
  "os"
  "os/exec"
  "os/user"
  "strconv"
  "syscall"
)

func isUserAuthenticated( ohmage_url string, username string, password string ) ( bool, error ) {
  if ohmage_url == "" || username == "" || password == "" {
    return false, errors.New( "Invalid/empty login parameters" )
  }
  ohmage := gohmage.NewClient( ohmage_url, "pam_ohmage" )
  auth_status, err := ohmage.UserAuthToken( username, password )
  if err != nil || auth_status == false {
    return false, errors.Wrap( err, "User authentication failed" )
  }
  user_info, err := ohmage.UserInfoRead( )
  if err != nil {
    return false, errors.Wrap( err, "Unable to fetch user information" )
  }
  user_classes, err := user_info.GetObject( "data", username, "classes" )
  if err != nil {
    return false, errors.Wrap( err, "Unable to find user classes" )
  }
  // todo: pick up test class from PAM configuration
  test_class, err := user_classes.GetString( "urn:class:rstudio" )
  if err != nil || test_class == "" {
    return false, errors.Wrap( err, "Test class not among the classes user is participant of" )
  }
  return true, nil
}

func isUserAccountReady( username string ) ( bool, error ) {
  // Check if the user account exists on the system & create if it doesn't
  user_account_id, err := getUserAccountId( username )
  if err != nil && user_account_id != 0 {
    return false, err
  }
  if user_account_id < 500 {
    return false, errors.New( "Cannot authenticate accounts with uid < 500" )
  } else if user_account_id == 0 {
    // todo: create the user account & home directory
  }

  user_home_directory_owner, user_home_directory_exists, err := getUserHomeDirectoryOwner( username )
  if err != nil {
    return false, err
  }
  if !user_home_directory_exists {
    // todo: create the user home directory
  }

  if user_home_directory_exists && user_home_directory_owner != user_account_id {
    // todo: chmod the home directory to this user
  }

  // Everything looks good! Account is ready to be used.
  return true, nil
}

func getUserAccountId( username string ) ( int, error ) {
  user_account, err := user.Lookup( username )
  if err != nil {
    return 0, errors.Wrap( err, "Unable to find a user" )
  }
  user_id, err := strconv.Atoi( user_account.Uid )
  if err != nil {
    return -1, err
  }
  return user_id, nil
}

func getUserHomeDirectoryOwner( username string ) ( int, bool, error ) {
  // sanitize username!
  directory, err := os.Open( "/home/" + username )
  defer directory.Close( )
  if err != nil {
    return 0, false, errors.Wrap( err, "Unable to read the home directory" )
  }
  directory_info, err := directory.Stat( )
  if err != nil {
    return 0, false, errors.Wrap( err, "Unable to read the home directory" )
  }
  if( !directory_info.IsDir( ) ) {
    return 0, false, errors.New( "Unable to find the home directory" )
  }
  directory_mode := directory_info.Sys( ).( *syscall.Stat_t )
  directory_owner := int( directory_mode.Uid )
  if directory_owner == 0 {
    return 0, false, errors.New( "Directory is owned by root. Cannot proceed" )
  }
  return directory_owner, true, nil
}