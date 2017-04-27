var gulp = require('gulp');

var spawn = require('child_process').spawn;

// How to define a new [platform]?
// Create ./test/platforms/[platform] directory (Check debian & centos for examples)
var platforms = [ 'centos', 'debian' ];

// steps when test_[PLATFORM] is called
// step 0
gulp.task( 'ohmage_test_setup', function( cb ) {
  spawnDockerCompose( [ 'up', 'ohmage_test_setup' ], cb );
} );

platforms.forEach( function( platform ) {

  gulp.task( 'build_' + platform, function( cb ) {
    spawnDockerCompose( [ '-f', './docker-compose.yml',
      '-f', './test/platforms/' + platform +'/docker-compose.yml',
      'build', '--no-cache', 'test_' + platform
    ], cb );
  } )

  // This is a good way to enforce the order in which the services are created.
  // ohmage_test_setup must be finished before starting webdriverio tests

  // step 1
  // this task could be directly called if you need to examine the containers
  // after the tests are run.
  gulp.task( 'start_test_' + platform, [ 'ohmage_test_setup' ], function ( cb ) {
    spawnDockerCompose( [ '-f', './docker-compose.yml',
      '-f', './test/platforms/' + platform +'/docker-compose.yml',
      'up', 'test_' + platform
    ], cb );
  } );

  // step 2
  gulp.task( 'test_' + platform, [ 'start_test_' + platform ], function ( cb ) {
    gulp.start( 'stop_test_' + platform );
  } );

  // step 3
  gulp.task( 'stop_test_' + platform, function ( cb ) {
    spawnDockerCompose( [ '-f', './docker-compose.yml',
      '-f', './test/platforms/' + platform + '/docker-compose.yml',
      'down'
    ] );
  } );
} );

function spawnDockerCompose( arguments, onExit ) {
  cmd = spawn( 'docker-compose', arguments );
  cmd.stdout.on( 'data', function (data) {
    process.stdout.write( data.toString( ) );
  });

  cmd.stderr.on( 'data', function (data) {
    process.stderr.write( data.toString( ) );
  });

  cmd.on( 'exit', function ( code ) {
    console.log( '[] docker-compose exited with code:' + code.toString( ) );
    !!onExit ? onExit( ) : null;
  });
}