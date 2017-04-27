import Ohmage from 'ohmage-es6';
import setup_config from './setup_config.json';
import util from 'util';

const ohmage_server = process.env.OHMAGE_SERVER || 'ohmage';
const ohmage_server_port = process.env.OHMAGE_SERVER_PORT || '8080';
let ohmage = new Ohmage( 'http://' + ohmage_server + ':' + ohmage_server_port + '/app', 'pam_ohmage_test' );

let auth_token = '';

let accounts = [ ].concat(  setup_config.login,
                            setup_config.directory_ownership,
                            setup_config.test_class_participation );
let accounts_to_add_to_class = [ ].concat(  setup_config.login,
                                            setup_config.directory_ownership );

setup_config.test_class_participation.forEach( account => {
  if( account.canSignIn ) {
    accounts_to_add_to_class.push( account )
  }
} );

accounts_to_add_to_class = accounts_to_add_to_class.map( account => account.username + ";restricted" );

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
      throw new Error( 'Unable to change admin password.' );
    }
    console.log( 'Authenticating admin' );
    return ohmage.authToken( 'ohmage.admin', 'ohmage.passwd' )
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
        throw error
      } );
  } )
  .then( ( ) => {
    console.log( 'Creating test class' );
    return ohmage.classCreate( setup_config.test_class_urn, "Test Class" )
            .catch( error => {
              if( error.parent.props.body.errors[ 0 ].code === '0900' ) {
                // class already exists error.. absorb this error!
                console.log( 'Test class already exists' );
                return true;
              } else {
                throw error;
              }
            } )
  } )
  .then( response => {
    if( ( typeof response === 'boolean' && response ) ||
        response.result === 'success' ) {
      console.log( 'Test class created' );
      return;
    } else {
      throw new Error( "Unable to create the test class" );
    }
  } )
  .then( ( ) => {
    console.log( 'Adding users to the test class' );
    return ohmage.classUpdate(  setup_config.test_class_urn,
                                {
                                  user_role_list_add: accounts_to_add_to_class.join( "," )
                                } )
  } )
  .then( response => {
    if( response.result === 'success' ) {
      console.log( 'Added users to the test class' );
      return;
    } else {
      throw new Error( "Unable to add users to the test class" );
    }
  } )
  .catch( error => {
    console.error( util.inspect( error, false, null ) );
  } )
