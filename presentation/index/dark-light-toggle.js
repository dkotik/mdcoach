// Embedded by the presentation index page as <dark-light-toggle>.

class DarkLightToggle extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onClick = this.onClick.bind(this)
    this.onKeyDown = this.onKeyDown.bind(this)
  }

  connectedCallback() {
    try {
      this.darkMode = window.localStorage.getItem('darkMode') !== 'false'
    } catch {
      this.darkMode = true
    }
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
    if ((event.code !== 'Enter' && event.code !== 'Space') || event.repeat) {
      return
    }

    event.preventDefault()
    this.toggle()
  }

  toggle() {
    this.darkMode = !this.darkMode
    try {
      window.localStorage.setItem('darkMode', String(this.darkMode))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
    this.applyTheme()
    this.updatePressedState()
  }
}

if (!customElements.get('dark-light-toggle')) {
  customElements.define('dark-light-toggle', DarkLightToggle)
}
