var fs = require( 'fs' );
var Launcher = require( 'webdriverio' ).Launcher;

var selenium_server = process.env.SELENIUM_SERVER || 'selenium';
var selenium_server_port = process.env.SELENIUM_SERVER_PORT || '4444';

var rstudio_server = process.env.RSTUDIO_SERVER || 'rstudio';
var rstudio_server_port = process.env.RSTUDIO_SERVER_PORT || '8787';

var wdio = new Launcher( './wdio.conf.js', {
  host: selenium_server,
  port: selenium_server_port,
  baseUrl: 'http://' + rstudio_server + ':' + rstudio_server_port
} );

console.log( 'Starting selenium test runner' )
wdio.run().then( function( code ) {
  var files = fs.readdirSync( './reports' );
  var json_report_file = files[ files.length - 1 ];
  var json_report = require( './reports/' + json_report_file );
  var state = json_report.state;
  if( state.failed > 0 ) {
    console.error( 'Some tests failed: ' + JSON.stringify( state ) );
  }
  process.exit( code );

}, function ( error ) {
    console.error( 'Launcher failed to start the test', error.stacktrace );
    process.exit( 1 );
} );