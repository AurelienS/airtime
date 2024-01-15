module.exports = {
    proxy: "localhost:3000", 
    port: 3001, // Port for BrowserSync
    files: "browser-refresh-trigger.nothing",
    open: false,
    reloadDelay: 50,
    injectChanges: true,
    watchOptions: {
        ignoreInitial: true
    }
};