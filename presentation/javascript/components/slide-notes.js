// use as custom element <slide-notes> to show notes from the focused slide

const notesPanelWidthStorageKey = 'presentationNotesWidth'
const notesPanelOpenStorageKey = 'presentationNotesOpen'
const notesPanelTabWidth = 40
const notesTextSizeStorageKey = 'mdcoachResizableTextGlobalSize'
const notesTextSizeMinimum = 20
const notesTextSizeMaximum = 180
const notesTextSizeDefault = 100
const notesTextSizeStep = 10
const showSlideNotesEventType = 'showSlideNotes'

class SlideNotes extends HTMLElement {
  static get observedAttributes() {
    return ['open']
  }

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onShowSlideNotes = this.onShowSlideNotes.bind(this)
    this.onTextSizeClick = this.onTextSizeClick.bind(this)
    this.onWindowResize = this.onWindowResize.bind(this)
    this.onResizeStart = this.onResizeStart.bind(this)
    this.onResizeMove = this.onResizeMove.bind(this)
    this.onResizeEnd = this.onResizeEnd.bind(this)
    this.onResizeKeyDown = this.onResizeKeyDown.bind(this)
    this.isResizing = false
    this.textSize = notesTextSizeDefault
  }

  connectedCallback() {
    if (!this.dock) {
      this.initializePanel()
    }

    this.textSizeToolbar.addEventListener('click', this.onTextSizeClick)
    this.resizeHandle.addEventListener('pointerdown', this.onResizeStart)
    this.resizeHandle.addEventListener('pointermove', this.onResizeMove)
    this.resizeHandle.addEventListener('pointerup', this.onResizeEnd)
    this.resizeHandle.addEventListener('pointercancel', this.onResizeEnd)
    this.resizeHandle.addEventListener('keydown', this.onResizeKeyDown)

    const defaultWidth = window.innerWidth * 0.36 + notesPanelTabWidth
    this.setPanelWidth(this.getStoredPanelWidth() ?? defaultWidth)
    window.addEventListener('resize', this.onWindowResize)
    window.addEventListener(showSlideNotesEventType, this.onShowSlideNotes)
    const storedOpenState = this.readOpenState()
    if (storedOpenState === null) {
      this.storeOpenState(this.hasAttribute('open'))
    } else {
      this.toggleAttribute('open', storedOpenState)
    }
    this.updateState()
  }

  createStyleElement() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        transform: translateX(calc(100% - 1.3rem));
        transition: transform 250ms ease;
        width: var(--notes-panel-width, calc(36vw + 1.3rem));
      }


      :host([open]) {
        transform: translateX(0);
      }

      .text-size-toolbar {
        float: right;
        font-size: 0.8rem;
        margin-bottom: 0.5rem;
      }

      .text-size-toolbar button {
        background: transparent;
        border: 0;
        color: inherit;
        cursor: pointer;
        font: inherit;
        line-height: 1;
        padding: 0.4rem 0.6rem;
      }

      .text-size-toolbar button:hover:not(:disabled) {
        border-color: var(--color-marker-background, #006eff);
      }

      .text-size-toolbar button:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      .text-size-toolbar button:disabled {
        cursor: default;
        opacity: 0.5;
      }

      .text-size-toolbar output {
        font-variant-numeric: tabular-nums;
        min-width: 4ch;
        text-align: center;
      }

      .notes-content {
        font-size: var(--notes-text-size, 100%);
      }

      .notes-content > :first-child {
        margin-top: 0;
        padding-top: 0;
      }

      .resize-handle:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      .resize-handle {
        bottom: 0;
        cursor: ew-resize;
        left: 1.25rem;
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

      @media (prefers-reduced-motion: reduce) {
        :host {
          transition: none;
        }
      }
    `
    return style
  }

  initializePanel() {
    const style = this.createStyleElement()

    this.dock = document.createElement('aside')
    this.dock.setAttribute('part', 'content')
    this.dock.setAttribute('aria-label', 'Speaker notes')
    this.content = document.createElement('div')
    this.textSizeToolbar = document.createElement('div')
    this.textSizeToolbar.className = 'text-size-toolbar'
    // this.textSizeToolbar.setAttribute('part', 'controls')
    this.textSizeToolbar.setAttribute('role', 'toolbar')
    this.decreaseTextSizeButton = this.createTextSizeButton(
      'Decrease font size',
      '－',
      -notesTextSizeStep,
    )
    this.increaseTextSizeButton = this.createTextSizeButton(
      'Increase font size',
      '＋',
      notesTextSizeStep,
    )
    this.textSizeOutput = document.createElement('output')
    this.textSizeOutput.setAttribute('aria-live', 'polite')
    this.textSizeOutput.setAttribute('aria-label', 'Current font size')
    this.textSizeToolbar.append(
      this.decreaseTextSizeButton,
      this.textSizeOutput,
      this.increaseTextSizeButton,
    )

    this.notesContent = document.createElement('div')
    this.notesContent.className = 'notes-content'
    this.notesContent.append(document.createElement('slot'))
    this.content.append(this.textSizeToolbar, this.notesContent)
    this.dock.append(this.content)

    this.loadTextSize()
    this.renderTextSize()


    this.resizeHandle = document.createElement('div')
    this.resizeHandle.className = 'resize-handle'
    this.resizeHandle.setAttribute('part', 'controls')
    this.resizeHandle.setAttribute('role', 'separator')
    this.resizeHandle.setAttribute('aria-label', 'Resize notes panel')
    this.resizeHandle.setAttribute('aria-orientation', 'vertical')
    this.resizeHandle.setAttribute('tabindex', '0')

    this.shadowRoot.replaceChildren(style, this.resizeHandle, this.dock)
  }

  createTextSizeButton(label, text, change) {
    const button = document.createElement('button')
    button.type = 'button'
    button.textContent = text
    button.setAttribute('aria-label', label)
    button.dataset.change = String(change)
    return button
  }

  onTextSizeClick(event) {
    const button = event.target.closest('button[data-change]')
    if (!button || !this.textSizeToolbar.contains(button) || button.disabled) {
      return
    }

    this.textSize = Math.min(
      notesTextSizeMaximum,
      Math.max(notesTextSizeMinimum, this.textSize + Number(button.dataset.change)),
    )
    this.renderTextSize()
    this.storeTextSize()
  }

  loadTextSize() {
    try {
      const storedSize = Number(window.localStorage.getItem(notesTextSizeStorageKey))
      this.textSize = Number.isFinite(storedSize) && storedSize >= notesTextSizeMinimum
        ? Math.min(storedSize, notesTextSizeMaximum)
        : notesTextSizeDefault
    } catch {
      this.textSize = notesTextSizeDefault
    }
  }

  renderTextSize() {
    this.notesContent.style.setProperty('--notes-text-size', `${this.textSize}%`)
    this.textSizeOutput.textContent = `${this.textSize}%`
    this.decreaseTextSizeButton.disabled = this.textSize <= notesTextSizeMinimum
    this.increaseTextSizeButton.disabled = this.textSize >= notesTextSizeMaximum
  }

  storeTextSize() {
    try {
      window.localStorage.setItem(notesTextSizeStorageKey, String(this.textSize))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

  disconnectedCallback() {
    this.textSizeToolbar?.removeEventListener('click', this.onTextSizeClick)
    window.removeEventListener('resize', this.onWindowResize)
    window.removeEventListener(showSlideNotesEventType, this.onShowSlideNotes)
    this.resizeHandle?.removeEventListener('pointerdown', this.onResizeStart)
    this.resizeHandle?.removeEventListener('pointermove', this.onResizeMove)
    this.resizeHandle?.removeEventListener('pointerup', this.onResizeEnd)
    this.resizeHandle?.removeEventListener('pointercancel', this.onResizeEnd)
    this.resizeHandle?.removeEventListener('keydown', this.onResizeKeyDown)
  }

  onShowSlideNotes(event) {
    if (typeof event.detail !== 'boolean') {
      return
    }

    this.toggleAttribute('open', event.detail)
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name !== 'open' || oldValue === newValue) {
      return
    }

    this.updateState()
    if (this.isConnected) {
      this.storeOpenState(newValue !== null)
    }
  }

  readOpenState() {
    try {
      const storedState = window.localStorage.getItem(notesPanelOpenStorageKey)
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
      window.localStorage.setItem(notesPanelOpenStorageKey, String(isOpen))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

  updateState() {
    const isOpen = this.hasAttribute('open')
    const toggleButton = document.querySelector('slide-notes-toggle')
    toggleButton?.setAttribute('aria-expanded', String(isOpen))
    toggleButton?.setAttribute('aria-label', isOpen ? 'Hide speaker notes' : 'Show speaker notes')
    toggleButton?.setAttribute('title', isOpen ? 'Hide notes' : 'Show notes')
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
}

if (!customElements.get('slide-notes')) {
  customElements.define('slide-notes', SlideNotes)
}
