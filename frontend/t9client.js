import * as config from './config.js'
import { fetchWithTimeout } from './fetch.js';

const digitsValidRegexp = /^[2-9]+$/; // only allow digits 2-9

const defaultBaseURL = `https://${config.server}/api/lookup/`
const defaultAjaxRequestTimeout = 6000 // expressed in ms

class t9Client {
    #baseURL
    #ajaxRequestTimeout

    constructor(config={}) {
        const {
            baseURL = defaultBaseURL,
            ajaxRequestTimeout = defaultAjaxRequestTimeout
        } = config;
        this.#baseURL = baseURL
        this.#ajaxRequestTimeout = ajaxRequestTimeout
    }

    getWords = (digits, exactOnly) => {
        return new Promise((resolve, reject) => {
            
            switch (true) {
            case digits.length === 0:
                resolve([]) // no digits, no words
                return
            case !digits.match(digitsValidRegexp):
                reject('Only the digits 2-9 are valid')
                return
            }
    
            const queryURL = `${this.#baseURL}${digits}?exact=${exactOnly}`

            fetchWithTimeout(queryURL, {
                method: 'GET',
                cache: 'force-cache',
                timeout: this.#ajaxRequestTimeout
            })
            .then((response) => response.json())
            .then(data => {
                if ( data.message && data.message.length > 0 ) {
                    reject(data.message)
                    return
                }
                resolve(data.words)
            })
            .catch((error) => {
                console.log(error)
                reject('There was an error communicating with the server')
            })
        })
    }
}

const defaultT9Client = new t9Client()
const getWords = defaultT9Client.getWords

export {
    t9Client,
    defaultT9Client,
    getWords,
}