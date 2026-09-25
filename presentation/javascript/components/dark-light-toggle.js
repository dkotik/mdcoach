// use as custom element <dark-light-toggle> anywhere on the page

class DarkLightToggle extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
  }

  connectedCallback() {
    this.darkMode = window.localStorage.getItem('darkMode') === 'true'
    this.applyTheme()

    const style = document.createElement('style')
    style.textContent = `
      button {
        align-items: center;
        display: inline-flex;
        gap: 0.4em;
      }

      svg {
        fill: currentColor;
        height: 1em;
        width: 1em;
      }
    `

    const button = document.createElement('button')
    button.type = 'button'
    button.append(this.createIcon(), document.createTextNode(this.label))
    button.addEventListener('click', () => this.toggle())
    this.shadowRoot.replaceChildren(style, button)
  }

  get label() {
    return this.darkMode ? 'Dark' : 'Light'
  }

  createIcon() {
    const icon = document.createElementNS(svgNamespace, 'svg')
    icon.setAttribute('viewBox', '0 0 24 24')
    icon.setAttribute('aria-hidden', 'true')
    const path = document.createElementNS(svgNamespace, 'path')
    path.setAttribute('d', 'M12 3a9 9 0 1 0 9 9c0-.5-.04-1-.12-1.48A7 7 0 0 1 13.48 3.12C13 3.04 12.5 3 12 3Z')
    icon.append(path)
    return icon
  }

  applyTheme() {
    document.documentElement.dataset.theme = this.darkMode ? 'dark' : 'light'
  }

  toggle() {
    this.darkMode = !this.darkMode
    if (this.darkMode) {
      window.localStorage.setItem('darkMode', 'true')
    } else {
      window.localStorage.removeItem('darkMode')
    }
    this.applyTheme()
    this.shadowRoot.querySelector('button').lastChild.textContent = this.label
  }
}

if (!customElements.get('dark-light-toggle')) {
  customElements.define('dark-light-toggle', DarkLightToggle)
}
