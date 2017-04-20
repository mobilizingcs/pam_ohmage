package main

import (
  "github.com/pkg/errors"
  "gohmage"
  "os"
  "os/exec"
  "os/user"
  "strconv"
  "syscall"
  "regexp"
  "github.com/op/go-logging"
)

var log = logging.MustGetLogger( "pam_ohmage" )

func init( ) {
  syslog, err := logging.NewSyslogBackend( "pam_ohmage" )
  if err != nil {
    return
  }
  leveledSyslog := logging.AddModuleLevel( syslog )
  format := logging.MustStringFormatter(
    `%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
  )
  formattedSyslog := logging.NewBackendFormatter( leveledSyslog, format )
  logging.SetBackend( formattedSyslog )
  logging.SetLevel( logging.CRITICAL, "pam_ohmage" )
}

func isUserAuthenticated( ohmage_url string, username string, password string ) ( bool, error ) {
  if ohmage_url == "" || username == "" || password == "" {
    log.Error( "Invalid/empty login parameters" )
    return false, errors.New( "Invalid/empty login parameters" )
  }
  ohmage := gohmage.NewClient( ohmage_url, "pam_ohmage" )
  auth_status, err := ohmage.UserAuthToken( username, password )
  if err != nil || auth_status == false {
    log.Error( "Ohmage authentication call failed" )
    return false, errors.Wrap( err, "User authentication failed" )
  }
  user_info, err := ohmage.UserInfoRead( )
  if err != nil {
    log.Error( "Ohmage user info call failed" )
    return false, errors.Wrap( err, "Unable to fetch user information" )
  }
  user_classes, err := user_info.GetObject( "data", username, "classes" )
  if err != nil {
    log.Error( "Failed to find user classes" )
    return false, errors.Wrap( err, "Unable to find user classes" )
  }
  // todo: pick up test class from PAM configuration
  test_class, err := user_classes.GetString( "urn:class:rstudio" )
  if err != nil || test_class == "" {
    log.Error( "Failed to find test class among user's classes" )
    return false, errors.Wrap( err, "Test class not among the classes user is participant of" )
  }
  return true, nil
}

func isUserAccountReady( username string ) ( bool, error ) {
  log.Debug( "called" )
  matched, err := regexp.MatchString( "([a-z_][a-z0-9_]{0,30})", username)
  if err != nil {
    log.Error( "Failed to validate username" )
    return false, err
  } else if !matched {
    log.Error( "Invalid username" )
    return false, errors.New( "Invalid username string" )
  }
  // Check if the user account exists on the system & create if it doesn't
  user_account_id, user_group_id, err := getUserAndGroupId( username )
  if err != nil && user_account_id == -1 && user_group_id == -1 {
    log.Error( "Unable to find uid and gid" )
    return false, err
  }
  if user_account_id == 0 {
    result, err := createUserAccount( username )
    if err != nil || result == false {
      log.Error( "Failed to create new user account" )
      return false, errors.Wrap( err, "Unable to create user account" )
    }
    user_account_id, user_group_id, err = getUserAndGroupId( username )
    if err != nil && user_account_id == -1 && user_group_id == -1 {
      log.Error( "Unable to find the newly created user's uid and gid" )
      return false, err
    }
  } else if user_account_id < 500 {
    log.Error( "Refusing to create new user account since uid < 500" )
    return false, errors.New( "Cannot authenticate accounts with uid < 500" )
  }

  user_home_directory_owner, user_home_directory_exists, err := getUserHomeDirectoryOwner( username )
  if err != nil {
    log.Error( "Could not check existence of user home directory" )
    return false, err
  }
  if !user_home_directory_exists {
    result, err := createUserHomeDirectory( username, user_account_id, user_group_id )
    if err != nil || result == false {
      log.Error( "Failed to create user home directory" )
      return false, errors.Wrap( err, "Unable to create user home directory" )
    }
  }

  if user_home_directory_exists && user_home_directory_owner != user_account_id {
    log.Debug( "User home directory exists but owner != uid" )
    result, err := setUserHomeDirectoryPermissions( username, user_account_id, user_group_id )
    if err != nil || result == false {
      log.Error( "Failed to set user home directory permissions" )
      return false, errors.Wrap( err, "Unable to set user home directory permissions" )
    }
  }

  log.Info( "User account is ready" )
  return true, nil
}

func getUserAndGroupId( username string ) ( int, int, error ) {
  log.Debug( "Looking up user" )
  user_account, err := user.Lookup( username )
  if err != nil {
    log.Debug( "Could not find user" )
    return 0, 0, errors.Wrap( err, "Unable to find a user" )
  }
  user_id, err := strconv.Atoi( user_account.Uid )
  if err != nil {
    log.Error( "Unable to parse uid" )
    return -1, -1, err
  }
  group_id, err := strconv.Atoi( user_account.Gid )
  if err != nil {
    log.Error( "Unable to parse gid" )
    return -1, -1, err
  }
  return user_id, group_id, nil
}

func getUserHomeDirectoryOwner( username string ) ( int, bool, error ) {
  log.Debug( "Looking up home directory owner" )
  directory, err := os.Open( "/home/" + username )
  defer directory.Close( )
  if err != nil {
    if os.IsNotExist( err ) {
      log.Debug( "Home directory does not exist" )
      return 0, false, nil
    } else {
      log.Error( "Unable to open home directory for reading" )
      return 0, false, errors.Wrap( err, "Unable to read the home directory" )
    }
  }
  directory_info, err := directory.Stat( )
  if err != nil {
    log.Error( "Unable to read home directory attributes" )
    return 0, false, errors.Wrap( err, "Unable to read the home directory" )
  }
  if( !directory_info.IsDir( ) ) {
    log.Error( "Home directory is not a directory (but a file)" )
    return 0, false, errors.New( "Unable to find the home directory" )
  }
  directory_mode := directory_info.Sys( ).( *syscall.Stat_t )
  directory_owner := int( directory_mode.Uid )
  if directory_owner == 0 {
    log.Error( "Home directory is owned by root. Some funny business is happening" )
    return 0, false, errors.New( "Directory is owned by root. Cannot proceed" )
  }
  return directory_owner, true, nil
}
func createUserAccount( username string ) ( bool, error ) {
  log.Debug( "Creating user account" )
  cmd := exec.Command(  "/usr/sbin/adduser",
                        "--shell",
                        "/usr/sbin/nologin",
                        "--disabled-password",
                        "--disabled-login",
                        "--no-create-home",
                        "-q",
                        "--gecos",
                        "\"\"",
                        username )
  _stdout_stderr, err := cmd.CombinedOutput( )
  if err != nil {
    log.Error( "Unable to obtain stderr & stdout" )
    return false, errors.Wrap( err, "Unable to execute the create users command")
  }
  // todo: check exit code of the command for more robustness
  stdout_stderr := string( _stdout_stderr )
  if stdout_stderr == "" {
    return true, nil
  } else {
    log.Error( "stderr for the command is non-empty: " + stdout_stderr )
    return false, errors.New( "Unable to execute the create users command. Output: " + stdout_stderr )
  }
}

func createUserHomeDirectory( username string, uid int, gid int ) ( bool, error ) {
  log.Debug( "Creating user home directory" )
  process_uid := syscall.Getuid( )
  process_gid := syscall.Getgid( )
  syscall.Setuid( uid )
  syscall.Setgid( gid )
  defer syscall.Setuid( process_uid )
  defer syscall.Setgid( process_gid )
  if err := os.Mkdir( "/home/" + username, 0750 ); err != nil {
    log.Error( "Failed to create user home directory" )
    return false, errors.Wrap( err, "Unable to create user home directory" )
  } else {
    perms, err := setUserHomeDirectoryPermissions( username, uid, gid )
    if err != nil {
      log.Error( "Failed to set user home directory ownership" )
      return false, errors.Wrap( err, "Unable to set the home directory ownership" )
    } else if perms {
      return true, nil
    }
  }
  return false, nil
}

func setUserHomeDirectoryPermissions( username string, uid int, gid int ) ( bool, error ) {
  log.Debug( "Setting user home directory ownership" )
  if err := os.Chown( "/home/" + username, uid, gid ); err != nil {
    log.Error( "Failed to set home directory ownership" )
    return false, err
  } else {
    return true, nil
  }
}