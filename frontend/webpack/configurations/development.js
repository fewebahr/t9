const environment = 'development'

function getEnvironment() {
    return environment
}

function getServer() {
    return '127.0.0.1:4239'
}

export {
    getEnvironment,
    getServer,
}