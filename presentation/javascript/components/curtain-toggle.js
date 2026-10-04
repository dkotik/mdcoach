// use as <curtain-toggle> to toggle the first presentation-curtain

class CurtainToggle extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onClick = this.onClick.bind(this)
    this.onKeyDown = this.onKeyDown.bind(this)
    this.onDOMReady = this.onDOMReady.bind(this)
    this.updateState = this.updateState.bind(this)
    this.curtainElement = null
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
    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'button')
    }
    if (!this.hasAttribute('tabindex')) {
      this.setAttribute('tabindex', '0')
    }

    this.addEventListener('click', this.onClick)
    this.addEventListener('keydown', this.onKeyDown)
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', this.onDOMReady, { once: true })
    } else {
      this.updateState()
    }
  }

  disconnectedCallback() {
    document.removeEventListener('DOMContentLoaded', this.onDOMReady)
    this.removeEventListener('click', this.onClick)
    this.curtainElement?.removeEventListener('change', this.updateState)
    this.curtainElement = null
  }

  onDOMReady() {
    this.updateState()
  }

  onClick() {
    const curtain = document.querySelector('presentation-curtain')
    if (!curtain) {
      return
    }

    if (curtain.IsDown()) {
      curtain.Up()
    } else {
      curtain.Down()
    }
  }

  onKeyDown(event) {
    if ((event.code !== 'Enter' && event.code !== 'Space') || event.repeat) {
      return
    }

    event.preventDefault()
    this.onClick()
  }

  updateState() {
    const curtain = document.querySelector('presentation-curtain')
    if (curtain !== this.curtainElement) {
      this.curtainElement?.removeEventListener('change', this.updateState)
      this.curtainElement = curtain
      this.curtainElement?.addEventListener('change', this.updateState)
    }

    const isOpen = Boolean(curtain?.IsDown())
    this.setAttribute('aria-pressed', String(isOpen))
    this.setAttribute('aria-label', isOpen ? 'Close presentation curtain' : 'Open presentation curtain')
  }
}

if (!customElements.get('curtain-toggle')) {
  customElements.define('curtain-toggle', CurtainToggle)
}
