var accounts = require( '../accounts.json' );

function signIn( username, password ) {
  $( '#username' ).setValue( username );
  $( '#password' ).setValue( password );
  $( '#buttonpanel > button > table > tbody > tr > td.inner' ).click( );
}

function userSignedIn( username ) {
  browser.waitUntil( function( ) {
    return  (
      $( 'div[title="' + username + '"]' ).isVisible( ) ||
      $( '#caption' ).isVisible( )
    )
  } );
  if( $( 'div[title="' + username + '"]' ).isVisible( ) ) {
    return true;
  } else if(  $( '#caption' ).isVisible( ) ) {
    return false;
  }
  return false;
}

function signOut( ) {
  browser.deleteCookie( );
  $( 'button[title="Sign out"]' ).click( );
  $( '#caption' ).waitForVisible( );
}

describe( 'RStudio', function() {

  beforeEach( function( ) {
    browser.url( '/' );
    $( '#caption' ).waitForVisible( );
  } );

  accounts.forEach( function( account ) {
    it( 'user '
      + account.username + ' should'
      + (!account.canSignIn ? ' not ' : ' ')
      + 'be able to sign in', function( ) {
      signIn( account.username, account.username )
      var user_signed_in = userSignedIn( account.username );
      if( account.canSignIn && user_signed_in ) {
        signOut( );
      } else if( !account.canSignIn && !user_signed_in ) {
        // cannot sign in.. and wasn't able to sign in either.. all good!
      } else {
        if( user_signed_in ) {
          signOut( );
        } else {
          // clear cookies just in case RStudio shows "Permission denied" page
          // with no controls
          browser.deleteCookie( );
        }
        throw new Error( "Sign-in test failed."
          + " Username: " + account.username
          + " canSignIn: " + account.canSignIn
          + " user_signed_in: " + user_signed_in );
      }
    } );
  } );

});