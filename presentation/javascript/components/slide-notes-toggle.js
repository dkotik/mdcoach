class SlideNotesToggle extends HTMLElement {
  static get observedAttributes() {
    return ['aria-expanded', 'aria-label', 'title']
  }

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onDOMReady = this.onDOMReady.bind(this)
    this.onClick = this.onClick.bind(this)
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.button = null
  }

  connectedCallback() {
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', this.onDOMReady, { once: true })
      return
    }

    this.initialize()
  }

  disconnectedCallback() {
    document.removeEventListener('DOMContentLoaded', this.onDOMReady)
    this.button?.removeEventListener('click', this.onClick)
    window.removeEventListener(navigationCompleteEventType, this.onNavigationComplete)
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue === newValue || !this.button) {
      return
    }

    if (newValue === null) {
      this.button.removeAttribute(name)
    } else {
      this.button.setAttribute(name, newValue)
    }
  }

  onDOMReady() {
    this.initialize()
  }

  initialize() {
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
          background: var(--color-menu-background);
          border: 1px solid var(--color-body-subtext);
          border-radius: 0.5rem 0 0 0.5rem;
          border-right: 0;
          color: var(--color-body-subtext, #aaa);
          cursor: pointer;
          font: inherit;
          font-size: 0.8rem;
          max-width: 2rem;
          overflow: hidden;
          padding: 0.75rem 0.5rem 0.75rem 0.1rem;
          pointer-events: auto;
          position: absolute;
          right: 0;
          top: 10vh;
          writing-mode: vertical-rl;
          z-index: 10;
        }

        button.notes-toggle[aria-expanded="true"] {
          border: 1px solid transparent;
          color: var(--color-body-subtext, #aaa);
          padding: 0.75rem 0.2rem 0.75rem 0.1rem;
        }

        button.notes-toggle:focus-visible {
          outline: 2px solid var(--color-marker-background, #006eff);
          outline-offset: 2px;
        }
      `

      this.button = document.createElement('button')
      this.button.id = 'slideNotesToggle'
      this.button.className = 'notes-toggle'
      this.button.type = 'button'
      this.button.setAttribute('aria-expanded', this.getAttribute('aria-expanded') ?? 'false')
      this.button.setAttribute('aria-label', this.getAttribute('aria-label') ?? 'Show speaker notes')
      this.button.setAttribute('title', this.getAttribute('title') ?? 'Show speaker notes')
      this.shadowRoot.replaceChildren(style, this.button)
    }

    this.button.addEventListener('click', this.onClick)
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    this.updateLabel(document.querySelector('main > section.is-focused'))
  }

  click() {
    if (!this.button || this.button.hidden) {
      return
    }

    this.button.click()
  }

  onClick() {
    if (this.button.hidden) {
      return
    }

    this.button.dispatchEvent(new CustomEvent('showSlideNotes', {
      bubbles: true,
      composed: true,
      detail: this.button.getAttribute('aria-expanded') !== 'true',
    }))
  }

  onNavigationComplete(event) {
    const slides = document.querySelectorAll('main > section')
    const slideIndex = Number(event.detail?.slideIndex)
    const section = Number.isInteger(slideIndex) && slideIndex >= 0
      ? slides[slideIndex]
      : document.querySelector('main > section.is-focused')
    this.updateLabel(section)
  }

  updateLabel(section) {
    const slideNotes = section?.querySelector('slide-notes')
    this.button.hidden = !slideNotes
    if (slideNotes) {
      this.button.textContent = slideNotes.getAttribute('label') ?? 'Notes'
    }
  }
}

if (!customElements.get('slide-notes-toggle')) {
  customElements.define('slide-notes-toggle', SlideNotesToggle)
}
