const isDev = import.meta.env.DEV

function getServer() {
    if ( isDev ) {
        return '127.0.0.1:4239'
    }
    // all other environments
    const host = location.hostname
    const port = location.port
    let server = host
    if ( port ) {
        server += ':' + port
    }
    return server
}

const server = getServer()

export {
    server,
    isDev
}