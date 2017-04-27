package main

import (
  "github.com/pkg/errors"
  "path/filepath"
  "os"
  "os/exec"
  "os/user"
  "strconv"
  "syscall"
  "regexp"
)

func isUserAccountReady( username string ) ( bool, error ) {
  log.Notice( "session module (open_session) called" )
  matched, err := regexp.MatchString( "([a-z_][a-z0-9_]{0,30})", username)
  if err != nil {
    log.Debug( "Failed to validate username" )
    return false, err
  } else if !matched {
    log.Debug( "Invalid username string" )
    return false, errors.New( "Invalid username string" )
  }
  // Check if the user account exists on the system & create if it doesn't
  user_account_id, user_group_id, err := getUserAndGroupId( username )
  if err != nil && user_account_id == -1 && user_group_id == -1 {
    log.Debug( "Unable to find uid and gid" )
    return false, err
  }
  if user_account_id == 0 {
    result, err := createUserAccount( username )
    if err != nil || result == false {
      log.Debug( "Failed to create new user account" )
      return false, errors.Wrap( err, "Unable to create user account" )
    }
    user_account_id, user_group_id, err = getUserAndGroupId( username )
    if err != nil && user_account_id == -1 && user_group_id == -1 {
      log.Debug( "Unable to find the newly created user's uid and gid" )
      return false, err
    }
  } else if user_account_id < 500 {
    log.Debug( "Refusing to create new user account since uid < 500" )
    return false, errors.New( "Cannot authenticate accounts with uid < 500" )
  }

  user_home_directory_owner, user_home_directory_exists, err := getUserHomeDirectoryOwner( username )
  if err != nil {
    log.Debug( "Could not check existence of user home directory" )
    return false, err
  }
  if !user_home_directory_exists {
    result, err := createUserHomeDirectory( username, user_account_id, user_group_id )
    if err != nil || result == false {
      log.Debug( "Failed to create user home directory" )
      return false, errors.Wrap( err, "Unable to create user home directory" )
    }
  }

  if user_home_directory_exists && user_home_directory_owner != user_account_id {
    log.Notice( "User home directory exists but owner != uid" )
    result, err := setUserHomeDirectoryPermissions( username, user_account_id, user_group_id )
    if err != nil || result == false {
      log.Debug( "Failed to set user home directory permissions" )
      return false, errors.Wrap( err, "Unable to set user home directory permissions" )
    }
  }

  log.Info( "User account is ready" )
  return true, nil
}

func getUserAndGroupId( username string ) ( int, int, error ) {
  log.Notice( "Looking up user" )
  user_account, err := user.Lookup( username )
  if err != nil {
    log.Debug( "Could not find user" )
    return 0, 0, errors.Wrap( err, "Unable to find a user" )
  }
  user_id, err := strconv.Atoi( user_account.Uid )
  if err != nil {
    log.Debug( "Unable to parse uid" )
    return -1, -1, err
  }
  group_id, err := strconv.Atoi( user_account.Gid )
  if err != nil {
    log.Debug( "Unable to parse gid" )
    return -1, -1, err
  }
  return user_id, group_id, nil
}

func getUserHomeDirectoryOwner( username string ) ( int, bool, error ) {
  log.Notice( "Looking up home directory owner" )
  directory, err := os.Open( "/home/" + username )
  defer directory.Close( )
  if err != nil {
    if os.IsNotExist( err ) {
      log.Notice( "Home directory does not exist" )
      return 0, false, nil
    } else {
      log.Debug( "Unable to open home directory for reading" )
      return 0, false, errors.Wrap( err, "Unable to read the home directory" )
    }
  }
  directory_info, err := directory.Stat( )
  if err != nil {
    log.Debug( "Unable to read home directory attributes" )
    return 0, false, errors.Wrap( err, "Unable to read the home directory" )
  }
  if( !directory_info.IsDir( ) ) {
    log.Debug( "Home directory is not a directory (but a file)" )
    return 0, false, errors.New( "Unable to find the home directory" )
  }
  directory_mode := directory_info.Sys( ).( *syscall.Stat_t )
  directory_owner := int( directory_mode.Uid )
  if directory_owner < 500 {
    log.Debug( "Home directory is owned by uid < 500. Cannot proceed." )
    return 0, false, errors.New( "Directory is owned by user uid < 500. Cannot proceed" )
  }
  return directory_owner, true, nil
}
func createUserAccount( username string ) ( bool, error ) {
  log.Notice( "Creating user account" )
  cmd := exec.Command(  "/usr/sbin/useradd",
                        "--shell",
                        "/usr/sbin/nologin",
                        "--no-create-home",
                        username )
  _stdout_stderr, err := cmd.CombinedOutput( )
  if err != nil {
    log.Debug( "Unable to obtain stderr & stdout" )
    return false, errors.Wrap( err, "Unable to execute the create users command")
  }
  // todo: check exit code of the command for more robustness
  stdout_stderr := string( _stdout_stderr )
  if stdout_stderr == "" {
    return true, nil
  } else {
    log.Debug( "stderr for the command is non-empty: " + stdout_stderr )
    return false, errors.New( "Unable to execute the create users command. Output: " + stdout_stderr )
  }
}

func createUserHomeDirectory( username string, uid int, gid int ) ( bool, error ) {
  log.Notice( "Creating user home directory" )
  process_uid := syscall.Getuid( )
  process_gid := syscall.Getgid( )
  syscall.Setuid( uid )
  syscall.Setgid( gid )
  defer syscall.Setuid( process_uid )
  defer syscall.Setgid( process_gid )
  if err := os.Mkdir( "/home/" + username, 0750 ); err != nil {
    log.Debug( "Failed to create user home directory" )
    return false, errors.Wrap( err, "Unable to create user home directory" )
  } else {
    perms, err := setUserHomeDirectoryPermissions( username, uid, gid )
    if err != nil {
      log.Debug( "Failed to set user home directory ownership" )
      return false, errors.Wrap( err, "Unable to set the home directory ownership" )
    } else if perms {
      return true, nil
    }
  }
  return false, nil
}

func ChownR(path string, uid int, gid int) error {
  return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
    if err == nil {
      err = os.Chown(name, uid, gid)
    }
    return err
  })
}

func setUserHomeDirectoryPermissions( username string, uid int, gid int ) ( bool, error ) {
  // todo: write webdriverio test to check accessibility from RStudio terminal
  log.Notice( "Setting user home directory ownership" )
  if err := ChownR( "/home/" + username, uid, gid ); err != nil {
    log.Debug( "Failed to set home directory ownership" )
    return false, err
  } else {
    return true, nil
  }
}