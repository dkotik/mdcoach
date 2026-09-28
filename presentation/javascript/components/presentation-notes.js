// use as custom element <presentation-notes> to show notes from the focused slide

const notesPanelWidthStorageKey = 'presentationNotesWidth'
const notesPanelTabWidth = 40

class PresentationNotes extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onToggle = this.onToggle.bind(this)
    this.onDOMReady = this.onDOMReady.bind(this)
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onWindowResize = this.onWindowResize.bind(this)
    this.onResizeStart = this.onResizeStart.bind(this)
    this.onResizeMove = this.onResizeMove.bind(this)
    this.onResizeEnd = this.onResizeEnd.bind(this)
    this.onResizeKeyDown = this.onResizeKeyDown.bind(this)
    this.isResizing = false
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        color: var(--color-menu-text, #666);
        display: block;
        height: 100vh;
        max-width: 80vw;
        position: fixed;
        right: 0;
        top: 0;
        transform: translateX(calc(100% - 2.5rem));
        transition: transform 250ms ease;
        width: var(--notes-panel-width, calc(36vw + 2.5rem));
        z-index: 9999;

        .dock {
          background: var(--color-menu-background, #eee);
          box-shadow: 0 0 1rem rgb(0 0 0 / 20%);
          box-sizing: border-box;
          height: 100%;
          margin-left: 2.5rem;
          overflow: auto;
          padding: 1rem 1.25rem;
        }
      }

      :host([open]) {
        transform: translateX(0);
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
        top: 10vh;
        writing-mode: vertical-rl;
        max-width: 4rem;
        overflow: hidden;
        z-index: 10;
      }

      button:focus-visible,
      .resize-handle:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      .resize-handle {
        bottom: 0;
        cursor: ew-resize;
        left: calc(2.5rem - 0.25rem);
        position: absolute;
        top: 0;
        touch-action: none;
        width: 0.5rem;
        z-index: 2;
      }

      .resize-handle::after {
        background: var(--color-body-subtext, #aaa);
        border-radius: 1px;
        content: '';
        height: 100%;
        left: calc(50% - 1px);
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        width: 2px;
      }

      .resize-handle:hover::after,
      .resize-handle:focus-visible::after {
        background: var(--color-marker-background, #006eff);
      }

      :host(:not([open])) .resize-handle {
        pointer-events: none;
      }

      :host([resizing]) {
        transition: none;
        user-select: none;
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

    this.dock = document.createElement('aside')
    this.dock.className = 'dock'
    this.dock.setAttribute('aria-label', 'Speaker notes')
    this.content = document.createElement('div')
    this.content.className = 'content'
    this.dock.append(this.content)

    this.resizeHandle = document.createElement('div')
    this.resizeHandle.className = 'resize-handle'
    this.resizeHandle.setAttribute('role', 'separator')
    this.resizeHandle.setAttribute('aria-label', 'Resize notes panel')
    this.resizeHandle.setAttribute('aria-orientation', 'vertical')
    this.resizeHandle.setAttribute('tabindex', '0')
    this.resizeHandle.addEventListener('pointerdown', this.onResizeStart)
    this.resizeHandle.addEventListener('pointermove', this.onResizeMove)
    this.resizeHandle.addEventListener('pointerup', this.onResizeEnd)
    this.resizeHandle.addEventListener('pointercancel', this.onResizeEnd)
    this.resizeHandle.addEventListener('keydown', this.onResizeKeyDown)

    this.shadowRoot.replaceChildren(style, this.button, this.resizeHandle, this.dock)
    const defaultWidth = window.innerWidth * 0.36 + notesPanelTabWidth
    this.setPanelWidth(this.getStoredPanelWidth() ?? defaultWidth)
    window.addEventListener('resize', this.onWindowResize)
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
    window.removeEventListener('resize', this.onWindowResize)
    this.resizeHandle?.removeEventListener('pointerdown', this.onResizeStart)
    this.resizeHandle?.removeEventListener('pointermove', this.onResizeMove)
    this.resizeHandle?.removeEventListener('pointerup', this.onResizeEnd)
    this.resizeHandle?.removeEventListener('pointercancel', this.onResizeEnd)
    this.resizeHandle?.removeEventListener('keydown', this.onResizeKeyDown)
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
    if (this.resizeHandle) {
      this.resizeHandle.inert = !isOpen
      this.resizeHandle.setAttribute('aria-hidden', String(!isOpen))
    }
  }

  onWindowResize() {
    this.setPanelWidth(this.getBoundingClientRect().width)
  }

  onResizeStart(event) {
    if (event.button !== 0) {
      return
    }

    event.preventDefault()
    this.isResizing = true
    this.resizeStartX = event.clientX
    this.resizeStartWidth = this.getBoundingClientRect().width
    this.resizePointerId = event.pointerId
    this.setAttribute('resizing', '')
    this.resizeHandle.setPointerCapture(event.pointerId)
  }

  onResizeMove(event) {
    if (!this.isResizing || event.pointerId !== this.resizePointerId) {
      return
    }

    this.setPanelWidth(this.resizeStartWidth + this.resizeStartX - event.clientX)
  }

  onResizeEnd(event) {
    if (!this.isResizing || event.pointerId !== this.resizePointerId) {
      return
    }

    this.isResizing = false
    this.resizePointerId = undefined
    this.removeAttribute('resizing')
    this.storePanelWidth(this.getBoundingClientRect().width)
  }

  onResizeKeyDown(event) {
    if (event.key !== 'ArrowLeft' && event.key !== 'ArrowRight') {
      return
    }

    event.preventDefault()
    const adjustment = event.key === 'ArrowLeft' ? 20 : -20
    this.setPanelWidth(this.getBoundingClientRect().width + adjustment, true)
  }

  setPanelWidth(width, persist = false) {
    const maxWidth = window.innerWidth * 0.8
    const minWidth = Math.min(240, maxWidth)
    const boundedWidth = Math.min(maxWidth, Math.max(minWidth, width))
    this.style.setProperty('--notes-panel-width', `${boundedWidth}px`)
    this.resizeHandle?.setAttribute('aria-valuemin', String(Math.round(minWidth)))
    this.resizeHandle?.setAttribute('aria-valuemax', String(Math.round(maxWidth)))
    this.resizeHandle?.setAttribute('aria-valuenow', String(Math.round(boundedWidth)))
    if (persist) {
      this.storePanelWidth(boundedWidth)
    }
  }

  getStoredPanelWidth() {
    try {
      const value = Number(window.localStorage.getItem(notesPanelWidthStorageKey))
      return Number.isFinite(value) && value > 0 ? value : null
    } catch {
      return null
    }
  }

  storePanelWidth(width) {
    try {
      window.localStorage.setItem(notesPanelWidthStorageKey, String(width))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

  updateNotes() {
    if (!this.content) {
      return
    }

    const slide = document.querySelector('main > aside.is-focused')
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
