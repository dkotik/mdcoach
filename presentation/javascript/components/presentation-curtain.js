// use as custom element <presentation-curtain> for a full-screen pause cover

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
    const style = document.createElement('style')
    style.textContent = `
      :host {
        display: block;
        inset: 0;
        pointer-events: none;
        position: fixed;
        z-index: 9999998;
      }

      :host([open]) {
        pointer-events: auto;
      }

      .cover {
        align-items: center;
        background: var(--color-body-background);
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

    this.updateState()
    window.addEventListener('keyup', this.onKeyUp)
  }

  disconnectedCallback() {
    window.removeEventListener('keyup', this.onKeyUp)
  }

  attributeChangedCallback() {
    this.updateState()
  }

  onKeyUp(event) {
    if (event.code !== 'Period') {
      return
    }

    event.preventDefault()
    this.toggleAttribute('open')
  }

  updateState() {
    const isOpen = this.hasAttribute('open')
    this.setAttribute('aria-hidden', String(!isOpen))
    this.inert = !isOpen
  }
}

if (!customElements.get('presentation-curtain')) {
  customElements.define('presentation-curtain', PresentationCurtain)
}
