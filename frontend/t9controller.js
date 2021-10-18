import { debounce } from './debounce.js';
import { defaultT9Client } from './t9client.js';

const defaultMinimumDigits = 2 // query will not be sent to server unless it has at least this many digits
const defaultSendRequestDelay = 150 // process user events only after the user has stopped typing for this amount of time (expressed in ms)
const defaultWordsDisplayLimit = 500 // user receives warning if query returns more than this number of words

export class t9Controller {
    #digitsInput
    #prefixDropdown
    #messageDiv
    #resultsDiv

    #t9Client

    #minimumDigits
    #wordsDisplayLimit

    #debouncedSendRequest

    constructor(
        digitsInput,
        prefixDropdown,
        messageDiv,
        resultsDiv,
        t9Client = defaultT9Client,
        options={},
    ) {
        const {
            minimumDigits = defaultMinimumDigits,
            sendEventDelay = defaultSendRequestDelay,
            wordsDisplayLimit = defaultWordsDisplayLimit
        } = options

        this.#digitsInput = digitsInput
        this.#prefixDropdown = prefixDropdown
        this.#messageDiv = messageDiv
        this.#resultsDiv = resultsDiv
        
        this.#t9Client = t9Client

        this.#minimumDigits = minimumDigits
        this.#wordsDisplayLimit = wordsDisplayLimit
        
        this.#debouncedSendRequest = debounce(this.#sendRequest, sendEventDelay)
    }

    start = () => {
        this.#digitsInput.addEventListener('keyup', this.#userEvent)
        this.#prefixDropdown.addEventListener('change', this.#userEvent)
    }

    #showMessage = (message, isError) => {
        this.#hideResults()
        this.#messageDiv.innerHTML = `<p>${message}</p>`
        if (isError) {
            this.#messageDiv.classList.add('error')
        } else {
            this.#messageDiv.classList.remove('error')
        }
        this.#messageDiv.classList.remove('hidden')
    }

    #hideMessage = () => {
        this.#messageDiv.classList.add('hidden')
    }

    #showResults = (words) => {
        this.#hideMessage()
        this.#resultsDiv.innerHTML = this.#getWordsHTML(words)
        this.#resultsDiv.classList.remove('hidden')
    }

    #hideResults = () => {
        this.#resultsDiv.classList.add('hidden')
    }

    #getWordsHTML = (words) => {
        words = words ? words : []
        let wordsHTML = ''
        words.forEach(word => {
            wordsHTML += '<div class="button disabled">'+ word +'</div>'
        })
        return wordsHTML
    }

    #refreshView = ({
        words = [],
        message = '',
        isError = false,
    }) => {
        if ( message ) {
            this.#showMessage(message, isError)
            return
        }
        if (words.length === 0) {
            this.#showMessage('No words matched your query')
            return
        }
        if (words.length > this.#wordsDisplayLimit) {
            this.#showMessage(`More than ${this.#wordsDisplayLimit} words matched your query`)
            return
        }

        this.#showResults(words)        
    }

    #userEvent = () => {
        console.log("processing event)")
        const digits = this.#digitsInput.value
        const exactOnly = this.#prefixDropdown.value !== 'prefix'

        if (digits.length === 0) {
            this.#hideMessage()
            return
        }
        if (digits.length < this.#minimumDigits) {
            this.#showMessage(`Enter at least ${this.#minimumDigits} digits`)
            return
        }
        this.#debouncedSendRequest(digits, exactOnly)
    }

    #sendRequest = (digits, exactOnly) => {
        this.#t9Client.getWords(digits, exactOnly)
        .then((words) => {
            this.#refreshView({
                words
            })
        })
        .catch((error) => {
            this.#refreshView({
                message:error,
                isError:true,
            })
        })
    }
}
