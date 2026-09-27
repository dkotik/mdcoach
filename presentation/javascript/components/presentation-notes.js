// use as custom element <presentation-notes> to show notes from the focused slide

class PresentationNotes extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onToggle = this.onToggle.bind(this)
    this.onDOMReady = this.onDOMReady.bind(this)
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        color: var(--color-menu-text, #666);
        display: block;
        height: 100vh;
        max-width: calc(100vw - 3rem);
        position: fixed;
        right: 0;
        top: 0;
        transform: translateX(calc(100% - 3rem));
        transition: transform 250ms ease;
        width: 24rem;
        z-index: 9999;
      }

      :host([open]) {
        transform: translateX(0);
      }

      .dock {
        background: var(--color-menu-background, #eee);
        box-shadow: 0 0 1rem rgb(0 0 0 / 20%);
        box-sizing: border-box;
        height: 100%;
        overflow: auto;
        padding: 1rem 1.25rem;
        padding-left: 4rem;
      }

      button {
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-body-subtext, #aaa);
        border-right: 0;
        border-radius: 0.5rem 0 0 0.5rem;
        color: inherit;
        cursor: pointer;
        font: inherit;
        font-weight: 600;
        left: 0;
        padding: 0.75rem 0.5rem;
        position: absolute;
        top: 1rem;
        writing-mode: vertical-rl;
      }

      button:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      h2 {
        font-size: 1.25rem;
        margin: 0 0 1rem;
      }

      .content:empty::after {
        color: var(--color-body-subtext, #aaa);
        content: 'No notes for this slide.';
      }

      @media (prefers-reduced-motion: reduce) {
        :host {
          transition: none;
        }
      }
    `

    this.button = document.createElement('button')
    this.button.type = 'button'
    this.button.textContent = 'Notes'
    this.button.addEventListener('click', this.onToggle)

    this.dock = document.createElement('section')
    this.dock.className = 'dock'
    this.dock.setAttribute('aria-label', 'Speaker notes')

    const heading = document.createElement('h2')
    heading.textContent = 'Notes'
    this.content = document.createElement('div')
    this.content.className = 'content'
    this.content.setAttribute('aria-live', 'polite')
    this.dock.append(heading, this.content)

    this.shadowRoot.replaceChildren(style, this.button, this.dock)
    this.updateState()
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', this.onDOMReady, { once: true })
    } else {
      this.updateNotes()
    }
    window.addEventListener('slideNavigationFinished', this.onNavigationComplete)
  }

  disconnectedCallback() {
    document.removeEventListener('DOMContentLoaded', this.onDOMReady)
    window.removeEventListener('slideNavigationFinished', this.onNavigationComplete)
    this.button?.removeEventListener('click', this.onToggle)
  }

  onToggle() {
    this.toggleAttribute('open')
    this.updateState()
  }

  onDOMReady() {
    this.updateNotes()
  }

  onNavigationComplete() {
    this.updateNotes()
  }

  updateState() {
    const isOpen = this.hasAttribute('open')
    this.button?.setAttribute('aria-expanded', String(isOpen))
    this.button?.setAttribute('aria-label', isOpen ? 'Hide speaker notes' : 'Show speaker notes')
    this.button?.setAttribute('title', isOpen ? 'Hide notes' : 'Show notes')
    if (this.dock) {
      this.dock.inert = !isOpen
      this.dock.setAttribute('aria-hidden', String(!isOpen))
    }
  }

  updateNotes() {
    if (!this.content) {
      return
    }

    const slide = document.querySelector('main > section.is-focused')
    const notes = slide?.querySelectorAll(':scope > .grid > .content > aside') || []
    const copies = Array.from(notes, (note) => {
      const copy = note.cloneNode(true)
      copy.removeAttribute('id')
      for (const element of copy.querySelectorAll('[id]')) {
        element.removeAttribute('id')
      }
      return copy
    })
    this.content.replaceChildren(...copies)
  }
}

if (!customElements.get('presentation-notes')) {
  customElements.define('presentation-notes', PresentationNotes)
}
