package main

/*
#cgo LDFLAGS: -lpam -fPIC
#define PAM_SM_AUTH
#define PAM_SM_SESSION
#include <security/pam_modules.h>
#include <stdlib.h>
#include <string.h>
char *string_from_argv(int i, char **argv) {
  return strdup(argv[i]);
}

char *get_username( pam_handle_t *pamHandle ) {
  if( !pamHandle ) {
    return NULL;
  }
  int pam_err = 0;
  const char *user;
  if ((pam_err = pam_get_item(pamHandle, PAM_USER, (const void**)&user)) != PAM_SUCCESS)
    return NULL;
  return strdup(user);
}

char *get_password( pam_handle_t *pamHandle ) {
  if( !pamHandle ) {
    return NULL;
  }
  struct pam_message _pamMessage;
  _pamMessage.msg_style = PAM_PROMPT_ECHO_OFF;
  _pamMessage.msg = "Please enter your ohmage password: ";

  struct pam_conv* pamConv;
  struct pam_response* pamResponse;
  const struct pam_message* pamMessage = &_pamMessage;

  const char *password = NULL;
  if( pam_get_item( pamHandle, PAM_CONV, (const void**)&pamConv ) == PAM_SUCCESS ) {
    if( pamConv ) {
      pamConv->conv( 1, &pamMessage, &pamResponse, pamConv->appdata_ptr );
      password = pamResponse[ 0 ].resp;
    }
  }

  return strdup( password );
}

PAM_EXTERN int pam_sm_authenticate(pam_handle_t* pamh, int flags, int argc, const char** argv) {
  return pamAuthenticate(pamh, flags, argc, argv);
}

PAM_EXTERN int pam_sm_setcred(pam_handle_t* pamh, int flags, int argc, const char** argv) {
  return PAM_IGNORE;
}

PAM_EXTERN int pam_sm_open_session(pam_handle_t* pamh, int flags, int argc, const char** argv) {
  return pamOpenSession(pamh, flags, argc, argv);
}

PAM_EXTERN int pam_sm_close_session(pam_handle_t* pamh, int flags, int argc, const char** argv) {
  return PAM_SUCCESS;
}
*/
import "C"
