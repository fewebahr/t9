/* loads configuration planted by webpack */
import * as config from 'config'

/* loads semantic dependencies */
import './semantic.js'

/* My own styles (one day) */
import './index.css'

/* third-party modules */
import lru from 'lru-cache'

const digitsValidRegexp = /^[2-9]+$/;
const baseUrl = 'https://'+ config.getServer() +'/api/lookup/'
const ajaxRequestTimeout = 6000 // expressed in ms
const sendRequestDelay = 50 // send request to server after user has stopped typing for this amount of time (expressed in ms)
const wordsDisplayLimit = 500 // user receives warning if query returns more than this number of words
const minimumDigits = 2 // query will not be sent to server unless it has at least this many digits
const cacheMaxBytes = 1024 * 1024 * 2 // 2 mb

$(document).ready(function() {

    const digitsInput = $('input#digits')
    const prefixDropdown = $('select.ui.dropdown')
    const errorDiv = $('div#error')
    const resultsDiv = $('div#results')
    const wordsDiv = $('div#words')

    // initialize lru cache
    let cache = lru({
        max: cacheMaxBytes,
        length: function(contents) { return contents.length }
    })

    const showError = (message) => {
        hideResults()
        errorDiv.text(message)
        errorDiv.removeClass('hidden')
    }

    const hideError = () => {
        errorDiv.addClass('hidden')
    }

    const showResults = (wordsHtml) => {
        hideError()
        wordsDiv.empty()
        wordsDiv.append(wordsHtml)
        resultsDiv.removeClass('hidden')
    }

    const hideResults = () => {
        resultsDiv.addClass('hidden')
    }

    const getWordsHtml = (words) => {
        words = words ? words : []

        let wordsHtml = ''

        switch (true) {
        case words.length === 0:
            wordsHtml = '<i class="exclamation triangle icon"></i> No words matched your query'
            break
        case words.length > wordsDisplayLimit:
            wordsHtml = '<i class="info circle icon"></i> More than '+ wordsDisplayLimit +' words matched your query.'
            break
        default:
            $.each(words, function( index, word ) {
                wordsHtml = wordsHtml.concat('<div class="ui blue basic button disabled">'+ word +'</div>')
            })
        }

        return wordsHtml
    }

    const sendRequest = () => {
        
        const digits = digitsInput.val()
        const exactOnly = prefixDropdown.val() === 'exact'
        if ( digits.length == 0 ) {
            return
        }

        const queryUrl = baseUrl + digits + (exactOnly ? '?exact=true' : '')
        
        const cacheKey = queryUrl
        if ( cache.has(cacheKey) ) {
            let wordsHtml = cache.get(cacheKey)
            showResults(wordsHtml)
            return
        }

        $.ajax({
            method: 'GET',
            url: queryUrl,
            cache: true,
            dataType: 'json',
            timeout: ajaxRequestTimeout,
        }).done(function(data) {
            hideError()
            if ( data.message && data.message.length > 0 ) {
                showError(data.message)
            } else {
                let wordsHtml = getWordsHtml(data.words)
                cache.set(cacheKey, wordsHtml)
                showResults(wordsHtml)
            }
        }).fail(function() {
            showError('There was an error communicating with the server')
        })
    }

    let doneTypingTimer = null
    const updateView = () => {

        clearTimeout(doneTypingTimer)

        const value = digitsInput.val()

        switch (true) {
        case value.length === 0:
            hideError()
            break
        case !value.match(digitsValidRegexp):
            showError('only the digits 2-9 are valid')
            break
        case value.length < minimumDigits:
            showError('enter at least '+ minimumDigits +' digits')
            break
        default:
            hideError()
            doneTypingTimer = setTimeout(sendRequest, sendRequestDelay)
            break
        }
    }

    
    prefixDropdown.dropdown(); // initialize the dropdown (required by semantic-ui)
    digitsInput.keyup(updateView)
    prefixDropdown.change(updateView)
})

