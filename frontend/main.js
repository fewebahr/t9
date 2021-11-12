import { ready } from './ready.js'
import { t9Controller } from './t9controller.js'
import { renderT9View } from './t9view.js'

ready(() => {
    renderT9View()
    new t9Controller(
        document.querySelector('input#digits'),
        document.querySelector('select#prefix'),
        document.querySelector('div#message'),
        document.querySelector('div#results'),
    ).start()
})

