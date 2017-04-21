package main

import (
  "github.com/op/go-logging"
  "github.com/pkg/errors"
  "net/url"
)

var log = logging.MustGetLogger( "pam_ohmage" )

func init( ) {
  syslog, err := logging.NewSyslogBackend( "pam_ohmage" )
  if err != nil {
    return
  }
  format := logging.MustStringFormatter(
    `%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
  )
  formattedSyslog := logging.NewBackendFormatter( syslog, format )
  leveledSyslog := logging.AddModuleLevel( formattedSyslog )
  logging.SetBackend( leveledSyslog )
  logging.SetLevel( logging.ERROR, "pam_ohmage" )
}

func parseOhmageUrl( ohmage_url string ) ( string, error ) {
  if ohmage_url == "" {
    return "", errors.New( "invalid ohmage url" )
  }
  parsed_ohmage_url, err := url.ParseRequestURI( ohmage_url )
  if err != nil {
    return "", errors.New( "invalid ohmage url" )
  } else {
    return parsed_ohmage_url.String( ), nil
  }
}