// use as <slideview-toggle> to toggle the body's slides class

const slideViewStorageKey = 'slideViewEnabled'

class SlideViewToggle extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onClick = this.onClick.bind(this)
    this.onKeyDown = this.onKeyDown.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        cursor: pointer;
        display: inline-flex;
      }

      :host(:focus-visible) {
        outline: 2px solid currentColor;
        outline-offset: 2px;
      }

      ::slotted(svg) {
        fill: currentColor;
        height: 1em;
        width: 1em;
      }
    `

    this.shadowRoot.replaceChildren(style, document.createElement('slot'))
    this.slideViewEnabled = this.readStoredState()

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'button')
    }
    if (!this.hasAttribute('tabindex')) {
      this.setAttribute('tabindex', '0')
    }

    this.addEventListener('click', this.onClick)
    this.addEventListener('keydown', this.onKeyDown)
    this.updateState()
  }

  disconnectedCallback() {
    this.removeEventListener('click', this.onClick)
    this.removeEventListener('keydown', this.onKeyDown)
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
    this.slideViewEnabled = !this.slideViewEnabled
    this.updateState()
    this.storeState()
  }

  updateState() {
    document.body?.classList.toggle('slides', this.slideViewEnabled)
    this.setAttribute('aria-pressed', String(this.slideViewEnabled))
    this.setAttribute(
      'aria-label',
      this.slideViewEnabled ? 'Disable slide view' : 'Enable slide view',
    )
  }

  readStoredState() {
    try {
      return window.localStorage.getItem(slideViewStorageKey) !== 'false'
    } catch {
      return true
    }
  }

  storeState() {
    try {
      window.localStorage.setItem(slideViewStorageKey, String(this.slideViewEnabled))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }
}

if (!customElements.get('slideview-toggle')) {
  customElements.define('slideview-toggle', SlideViewToggle)
}
