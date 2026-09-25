// use as custom element <presentation-menu> around presentation controls

class PresentationMenu extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onMouseMove = this.onMouseMove.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        display: block;
        left: 50%;
        position: fixed;
        top: 0;
        transform: translate(-50%, -100%);
        transition: transform 220ms ease;
        z-index: 1000;
      }

      :host([open]) {
        transform: translate(-50%, 0);
      }

      nav {
        align-items: center;
        background: var(--color-menu-background, #eee);
        border-radius: 0 0 0.5rem 0.5rem;
        box-shadow: 0 0.25rem 1rem rgb(0 0 0 / 20%);
        box-sizing: border-box;
        color: var(--color-menu-text, #666);
        display: flex;
        gap: 0.5rem;
        max-width: calc(100vw - 1rem);
        padding: 0.5rem;
      }

      @media (prefers-reduced-motion: reduce) {
        :host {
          transition: none;
        }
      }
    `

    this.navigation = document.createElement('nav')
    this.navigation.setAttribute('aria-label', this.getAttribute('label') || 'Presentation menu')
    const slot = document.createElement('slot')
    this.navigation.append(slot)
    this.shadowRoot.replaceChildren(style, this.navigation)

    this.updateVisibility(this.hasAttribute('open'))
    window.addEventListener('mousemove', this.onMouseMove, { passive: true })
  }

  disconnectedCallback() {
    window.removeEventListener('mousemove', this.onMouseMove)
  }

  onMouseMove(event) {
    const revealBoundary = window.innerHeight * 0.15
    const shouldOpen = event.clientY <= revealBoundary || this.matches(':hover')
    this.updateVisibility(shouldOpen)
  }

  updateVisibility(isVisible) {
    this.toggleAttribute('open', isVisible)
    if (!this.navigation) {
      return
    }

    this.navigation.inert = !isVisible
    this.navigation.setAttribute('aria-hidden', String(!isVisible))
  }
}

if (!customElements.get('presentation-menu')) {
  customElements.define('presentation-menu', PresentationMenu)
}
