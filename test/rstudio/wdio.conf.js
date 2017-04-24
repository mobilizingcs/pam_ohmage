exports.config = {
    host: 'selenium',
    port: 4444,
    path: '/wd/hub',
    specs: [
        './specs/**/*.js'
    ],
    reporters: ['dot', 'json'],
    reporterOptions: {
      outputDir: './reports'
    },
    exclude: [ ],
    maxInstances: 10,
    capabilities: [{
        maxInstances: 5,
        browserName: 'chrome'
    }],
    sync: true,
    // Level of logging verbosity: silent | verbose | command | data | result | error
    logLevel: 'error',
    coloredLogs: true,
    bail: 0,
    screenshotPath: './errorShots/',
    baseUrl: 'http://rstudio:8787',
    waitforTimeout: 20000,
    connectionRetryTimeout: 90000,
    connectionRetryCount: 3,
    framework: 'mocha',
    mochaOpts: {
        ui: 'bdd'
    },
}
