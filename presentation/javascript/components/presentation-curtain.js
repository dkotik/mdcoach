// use as custom element <presentation-curtain> for a full-screen pause cover

const curtainOpenStorageKey = 'presentationCurtainOpen'

class PresentationCurtain extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onKeyUp = this.onKeyUp.bind(this)
  }

  static get observedAttributes() {
    return ['open']
  }

  connectedCallback() {
    const storedOpenState = this.readOpenState()
    const style = document.createElement('style')
    style.textContent = `
      :host {
        display: block;
        inset: 0;
        pointer-events: none;
        position: fixed;
        z-index: 900;
      }

      :host([open]) {
        pointer-events: auto;
      }

      .cover {
        align-items: center;
        background: linear-gradient(
          135deg,
          color-mix(in srgb, var(--color-marker-background, #006eff) 8%, var(--color-body-background, white)) 0%,
          color-mix(in srgb, var(--color-marker-background, #006eff) 12%, var(--color-body-background, white)) 33%,
          color-mix(in srgb, var(--color-marker-background, #006eff) 16%, var(--color-body-background, white)) 66%,
          color-mix(in srgb, var(--color-marker-background, #006eff) 20%, var(--color-body-background, white)) 100%
        );
        border-bottom: 0.2rem solid var(--color-body-subtext, #aaa);
        box-sizing: border-box;
        color: var(--color-body-subtext, #aaa);
        display: grid;
        font-size: 3rem;
        font-weight: bold;
        height: 100%;
        justify-content: center;
        transform: translateY(-100%);
        transition: transform 1450ms ease;
        width: 100%;
      }

      :host([open]) .cover {
        transform: translateY(0);
      }

      @media (prefers-reduced-motion: reduce) {
        .cover {
          transition: none;
        }
      }
    `

    if (!this.cover) {
      this.cover = document.createElement('div')
      this.cover.className = 'cover'
      this.cover.textContent = this.getAttribute('label') || '⏸'
      this.cover.setAttribute('aria-hidden', 'true')
    }
    this.shadowRoot.replaceChildren(style, this.cover)

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'status')
    }
    if (!this.hasAttribute('aria-label')) {
      this.setAttribute('aria-label', 'Presentation paused')
    }

    if (storedOpenState === null) {
      this.storeOpenState(this.IsDown())
    } else if (storedOpenState) {
      this.Down()
    } else {
      this.Up()
    }
    this.updateState()
    window.addEventListener('keyup', this.onKeyUp)
  }

  disconnectedCallback() {
    window.removeEventListener('keyup', this.onKeyUp)
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue === newValue) {
      return
    }

    this.updateState()
    if (this.isConnected) {
      this.storeOpenState(this.IsDown())
    }
  }

  Down() {
    if (this.IsDown()) {
      return
    }
    this.setAttribute('open', '')
    this.dispatchEvent(new Event('change', { bubbles: true, composed: true }))
  }

  Up() {
    if (!this.IsDown()) {
      return
    }
    this.removeAttribute('open')
    this.dispatchEvent(new Event('change', { bubbles: true, composed: true }))
  }

  IsDown() {
    return this.hasAttribute('open')
  }

  onKeyUp(event) {
    if (event.code !== 'Period') {
      return
    }

    event.preventDefault()
    if (this.IsDown()) {
      this.Up()
    } else {
      this.Down()
    }
  }

  readOpenState() {
    try {
      const storedState = window.localStorage.getItem(curtainOpenStorageKey)
      if (storedState === 'true') {
        return true
      }
      if (storedState === 'false') {
        return false
      }
    } catch {
      return null
    }
    return null
  }

  storeOpenState(isOpen) {
    try {
      window.localStorage.setItem(curtainOpenStorageKey, String(isOpen))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

  updateState() {
    const isOpen = this.IsDown()
    this.setAttribute('aria-hidden', String(!isOpen))
    this.inert = !isOpen
  }
}

if (!customElements.get('presentation-curtain')) {
  customElements.define('presentation-curtain', PresentationCurtain)
}
