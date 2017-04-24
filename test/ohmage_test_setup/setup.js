import Ohmage from 'ohmage-es6';
import accounts from './accounts.json';

const ohmage_server = process.env.OHMAGE_SERVER || 'ohmage';
const ohmage_server_port = process.env.OHMAGE_SERVER_PORT || '8080';
let ohmage = new Ohmage( 'http://' + ohmage_server + ':' + ohmage_server_port + '/app', 'pam_ohmage_test' );

let auth_token = '';

function createAccountIfNeeded(  ) {
  let account = accounts.pop( );
  const username = account.username;
  return ohmage.userCreate( {
            admin: false,
            enabled: true,
            new_account: false,
            username: username,
            password: username
          } )
          .then( response => {
            console.log( username, 'created' );
            return createAllAccounts( );
          } )
          .catch( error => {
            // if the error is not related to user already existing, rethrow it
            if ( error.parent.props.body.errors[0].code !== '1000' ) {
              throw error;
            }
            console.log( username, 'already exists' );
            return createAllAccounts( );
          } );
}

function createAllAccounts( ) {
  if( accounts.length === 0 ) return true
  return createAccountIfNeeded( )
          .then( ( ) => {
            if( accounts.length > 0 ) {
              return createAccountIfNeeded( );
            } else {
              return true;
            }
          } );
}

console.log( 'Changing admin password' );
ohmage.changePassword( 'ohmage.admin', 'ohmage.passwd', 'ohmage.passwd' )
  .then( response => {
    console.log( 'Admin password changed' );
    if( response.result !== 'success' ) {
      return;
    }
    console.log( 'Authenticating admin' );
    ohmage.authToken( 'ohmage.admin', 'ohmage.passwd' )
      .then( response => {
        console.log( 'Authentication successful' );
        if( response.result === 'success' ) {
          auth_token = response.token;
          ohmage._setToken( auth_token );
        }
        console.log( 'Creating accounts:' )
      } )
      .then( createAllAccounts )
      .then( ( ) => {
        console.log( 'Account creation complete.' );
      } )
      .catch( error => {
        console.error( error );
      } );
  } )
  .catch( error => {
    console.error( error );
  } )
