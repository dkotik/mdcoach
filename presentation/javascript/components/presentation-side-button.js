// use as <presentation-side-button side="left|right"> around side controls

class PresentationSideButton extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onMouseMove = this.onMouseMove.bind(this)
  }

  connectedCallback() {
    if (this.getAttribute('side') !== 'right') {
      this.setAttribute('side', 'left')
    }

    const style = document.createElement('style')
    style.textContent = `
      :host {
        display: block;
        position: fixed;
        top: 50%;
        transition: transform 220ms ease;
        z-index: 1000;
        cursor: pointer;
      }

      :host([side="left"]) {
        left: 0;
        transform: translate(-100%, -50%);
      }

      :host([side="left"][open]) {
        transform: translate(0, -50%);
      }

      :host([side="right"]) {
        right: 0;
        transform: translate(100%, -50%);
      }

      :host([side="right"][open]) {
        transform: translate(0, -50%);
      }

      aside {
        align-items: center;
        background: var(--color-menu-background, #eee);
        box-shadow: 0 0.25rem 1rem rgb(0 0 0 / 20%);
        box-sizing: border-box;
        color: var(--color-menu-text, #666);
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        max-height: calc(100vh - 2rem);
        overflow: auto;
        padding: 0.5rem;
      }

      :host([side="left"]) aside {
        border-radius: 0 0.5rem 0.5rem 0;
      }

      :host([side="right"]) aside {
        border-radius: 0.5rem 0 0 0.5rem;
      }

      @media (prefers-reduced-motion: reduce) {
        :host {
          transition: none;
        }
      }
    `

    this.panel = document.createElement('aside')
    this.panel.setAttribute('aria-label', this.getAttribute('label') || 'Presentation side controls')
    this.panel.append(document.createElement('slot'))
    this.shadowRoot.replaceChildren(style, this.panel)

    this.updateVisibility(this.hasAttribute('open'))
    window.addEventListener('mousemove', this.onMouseMove, { passive: true })
  }

  disconnectedCallback() {
    window.removeEventListener('mousemove', this.onMouseMove)
  }

  onMouseMove(event) {
    const edgeDistance = window.innerWidth * 0.15
    const nearDock = this.getAttribute('side') === 'right'
      ? event.clientX >= window.innerWidth - edgeDistance
      : event.clientX <= edgeDistance
    const shouldOpen = nearDock || this.matches(':hover') || this.matches(':focus-within')
    this.updateVisibility(shouldOpen)
  }

  updateVisibility(isVisible) {
    this.toggleAttribute('open', isVisible)
    if (!this.panel) {
      return
    }

    this.panel.inert = !isVisible
    this.panel.setAttribute('aria-hidden', String(!isVisible))
  }
}

if (!customElements.get('presentation-side-button')) {
  customElements.define('presentation-side-button', PresentationSideButton)
}
