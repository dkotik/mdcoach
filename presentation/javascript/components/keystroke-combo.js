// use as <keystroke-combo timeout="1000">; timeout is in milliseconds

const keyStrokeCompleteEventType = "keystroke-combo-complete"

class KeystrokeCombo extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.keys = []
    this.recording = false
    this.timeout = undefined
    this.onClick = this.onClick.bind(this)
    this.onKeyDown = this.onKeyDown.bind(this)
    this.finishRecording = this.finishRecording.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        cursor: text;
        display: inline-flex;
        max-width: 100%;
      }

      :host(:focus-visible) output {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      output {
        align-items: center;
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-menu-text, #666);
        border-radius: 0.25rem;
        color: var(--color-menu-text, #666);
        display: flex;
        flex-wrap: wrap;
        gap: 0.35rem;
        min-height: 1.5rem;
        min-width: 6rem;
        padding: 0.25rem 0.5rem;
      }

      :host([recording]) output {
        border-color: var(--color-marker-background, #006eff);
      }

      kbd {
        background: var(--color-body-background, white);
        border: 1px solid currentColor;
        border-radius: 0.2rem;
        font: inherit;
        padding: 0.05rem 0.3rem;
      }

      .separator,
      .placeholder {
        opacity: 0.65;
      }
    `

    const pattern = this.getAttribute('filter')
    if (pattern === null) {
      this.matchesFilter = (key) => true
    } else {
      const filter = new RegExp(pattern)
      this.matchesFilter = (key) => filter.test(key)
    }

    this.output = document.createElement('output')
    this.output.setAttribute('aria-live', 'polite')
    this.shadowRoot.replaceChildren(style, this.output)

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'group')
    }
    if (!this.hasAttribute('aria-label')) {
      this.setAttribute('aria-label', 'Keystroke sequence recorder')
    }
    if (!this.hasAttribute('tabindex')) {
      this.setAttribute('tabindex', '0')
    }

    this.addEventListener('click', this.onClick)
    window.addEventListener('keydown', this.onKeyDown)
    this.renderKeys()
  }

  disconnectedCallback() {
    window.clearTimeout(this.timeout)
    this.timeout = undefined
    this.recording = false
    this.removeAttribute('recording')
    this.removeEventListener('click', this.onClick)
    window.removeEventListener('keydown', this.onKeyDown)
  }

  onClick() {
    this.focus({ preventScroll: true })
  }

  onKeyDown(event) {
    if (event.repeat || event.isComposing || !this.matchesFilter(event.key)) {
      return
    }

    const key = event.key === ' ' ? 'Space' : event.key
    if (!key) {
      return
    }

    event.preventDefault()
    event.stopPropagation()

    if (!this.recording) {
      this.keys = []
      this.recording = true
      this.setAttribute('recording', '')
    }

    this.keys.push(key)
    this.renderKeys()
    window.clearTimeout(this.timeout)
    this.timeout = window.setTimeout(
      this.finishRecording,
      this.getTimeoutDuration(),
    )
  }

  getTimeoutDuration() {
    const value = this.getAttribute('timeout')
    if (value === null) {
      return 1000
    }

    const duration = Number(value)
    return Number.isFinite(duration) && duration >= 0 ? duration : 1000
  }

  renderKeys() {
    if (this.keys.length === 0) {
      const placeholder = document.createElement('span')
      placeholder.className = 'placeholder'
      placeholder.textContent = this.getAttribute('placeholder') || 'Press keys'
      this.output.replaceChildren(placeholder)
      return
    }

    const contents = []
    for (const [index, key] of this.keys.entries()) {
      // if (index > 0) {
      //   const separator = document.createElement('span')
      //   separator.className = 'separator'
      //   separator.textContent = '→'
      //   contents.push(separator)
      // }
      const keycap = document.createElement('kbd')
      keycap.textContent = key
      contents.push(keycap)
    }
    this.output.replaceChildren(...contents)
  }

  finishRecording() {
    if (!this.recording) {
      return
    }

    this.recording = false
    this.timeout = undefined
    this.removeAttribute('recording')
    this.dispatchEvent(new CustomEvent(
      keyStrokeCompleteEventType, {
      bubbles: true,
      composed: true,
      detail: this.keys.join(' '),
    }))
  }
}

if (!customElements.get('keystroke-combo')) {
  customElements.define('keystroke-combo', KeystrokeCombo)
}
