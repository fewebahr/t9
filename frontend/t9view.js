import normalizeCSS from './normalize.css'
import skeletonCSS from './skeleton.css'
import appCSS from './style.css'

import html from './t9view.html?raw'

function getTopLevelElement(tagName) {
    const elements = document.getElementsByTagName(tagName);
    if (elements.length >= 1) {
        return elements[0];
    }

    const element = document.createElement(tagName);
    getHtmlElement().appendChild(element);
    return element;
}

function insertCSS(cssString) {
    const styleElement = document.createElement("style");
    styleElement.appendChild(document.createTextNode(cssString));
    getTopLevelElement("head").appendChild(styleElement);
}

function insertHTML(htmlString) {
    getTopLevelElement("body").innerHTML = html;
}

export function renderT9View() {
    insertCSS(normalizeCSS)
    insertCSS(skeletonCSS)
    insertCSS(appCSS)
    insertHTML(html)
}
