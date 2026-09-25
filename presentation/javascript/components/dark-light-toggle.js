// use as custom element <dark-light-toggle> around document-provided content

class DarkLightToggle extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onClick = this.onClick.bind(this)
    this.onKeyDown = this.onKeyDown.bind(this)
  }

  connectedCallback() {
    this.darkMode = window.localStorage.getItem('darkMode') === 'true'
    this.applyTheme()

    const style = document.createElement('style')
    style.textContent = `
      :host {
        align-items: center;
        cursor: pointer;
        display: inline-flex;
      }

      ::slotted(svg) {
        fill: currentColor;
        height: 1em;
        margin-inline-start: 0.4em;
        width: 1em;
      }

      :host(:focus-visible) {
        outline: 2px solid currentColor;
        outline-offset: 2px;
      }
    `

    this.shadowRoot.replaceChildren(style, document.createElement('slot'))

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'button')
    }
    if (!this.hasAttribute('tabindex')) {
      this.setAttribute('tabindex', '0')
    }

    this.addEventListener('click', this.onClick)
    this.addEventListener('keydown', this.onKeyDown)
    this.updatePressedState()
  }

  disconnectedCallback() {
    this.removeEventListener('click', this.onClick)
    this.removeEventListener('keydown', this.onKeyDown)
  }

  applyTheme() {
    document.documentElement.dataset.theme = this.darkMode ? 'dark' : 'light'
  }

  updatePressedState() {
    this.setAttribute('aria-pressed', String(this.darkMode))
  }

  onClick() {
    this.toggle()
  }

  onKeyDown(event) {
    if (event.code !== 'Enter' && event.code !== 'Space') {
      return
    }

    event.preventDefault()
    this.toggle()
  }

  toggle() {
    this.darkMode = !this.darkMode
    if (this.darkMode) {
      window.localStorage.setItem('darkMode', 'true')
    } else {
      window.localStorage.removeItem('darkMode')
    }
    this.applyTheme()
    this.updatePressedState()
  }
}

if (!customElements.get('dark-light-toggle')) {
  customElements.define('dark-light-toggle', DarkLightToggle)
}
