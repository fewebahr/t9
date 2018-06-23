const environment = 'production'

function getEnvironment() {
    return environment
}

function getServer() {
    return location.hostname+(location.port ? ':'+location.port: '')
}

export {
    getEnvironment,
    getServer,
}