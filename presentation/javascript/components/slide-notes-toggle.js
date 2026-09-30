class SlideNotesToggle extends HTMLElement {
  static get observedAttributes() {
    return ['label', 'open']
  }

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onClick = this.onClick.bind(this)
  }

  connectedCallback() {
    if (!this.button) {
      const style = document.createElement('style')
      style.textContent = `
        :host {
          display: block;
          inset: 0;
          pointer-events: none;
          position: absolute;
          z-index: 10;
        }

        button.notes-toggle {
          background: transparent;
          border: 1px solid var(--color-body-subtext, #aaa);
          border-radius: 0.5rem 0 0 0.5rem;
          border-right: 0;
          color: var(--color-body-subtext, #aaa);
          cursor: pointer;
          font: inherit;
          font-size: 0.8rem;
          left: 0;
          max-width: 2rem;
          overflow: hidden;
          padding: 0.75rem 0.5rem 0.75rem 0.1rem;
          pointer-events: auto;
          position: absolute;
          top: 10vh;
          writing-mode: vertical-rl;
          z-index: 10;
        }

        :host([open]) button.notes-toggle {
          border: 1px solid transparent;
          color: var(--color-body-subtext, #aaa);
          left: auto;
          padding: 0.75rem 0.2rem 0.75rem 0.1rem;
          right: 0;
        }

        button.notes-toggle:focus-visible {
          outline: 2px solid var(--color-marker-background, #006eff);
          outline-offset: 2px;
        }
      `

      this.button = document.createElement('button')
      this.button.className = 'notes-toggle'
      this.button.type = 'button'
      this.shadowRoot.replaceChildren(style, this.button)
    }

    this.button.addEventListener('click', this.onClick)
    this.updateState()
  }

  disconnectedCallback() {
    this.button?.removeEventListener('click', this.onClick)
  }

  attributeChangedCallback() {
    this.updateState()
  }

  onClick() {
    const root = this.getRootNode()
    const slideNotes = root.host ?? this.closest('slide-notes')
    if (!slideNotes || slideNotes.localName !== 'slide-notes') {
      return
    }

    slideNotes.toggleAttribute('open')
  }

  updateState() {
    if (!this.button) {
      return
    }

    const isOpen = this.hasAttribute('open')
    this.button.textContent = this.getAttribute('label') ?? 'Notes'
    this.button.setAttribute('aria-expanded', String(isOpen))
    this.button.setAttribute('aria-label', isOpen ? 'Hide speaker notes' : 'Show speaker notes')
    this.button.setAttribute('title', isOpen ? 'Hide notes' : 'Show notes')
  }
}

if (!customElements.get('slide-notes-toggle')) {
  customElements.define('slide-notes-toggle', SlideNotesToggle)
}
