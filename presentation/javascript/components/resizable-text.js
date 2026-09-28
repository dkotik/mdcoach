// use as <resizable-text data-key="slide-text-size">; size persists as a percentage

const resizableTextMinimum = 20
const resizableTextMaximum = 180
const resizableTextDefault = 100
const resizableTextStep = 10

class ResizableText extends HTMLElement {
  static get observedAttributes() {
    return ['data-key']
  }

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.size = resizableTextDefault
    this.onClick = this.onClick.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        color: var(--color-menu-text, #666);
        display: block;
        max-width: 100%;
      }

      .toolbar {
        align-items: center;
        display: flex;
        gap: 0.4rem;
        margin-bottom: 0.5rem;
      }

      button {
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-body-subtext, #aaa);
        border-radius: 0.3rem;
        color: inherit;
        cursor: pointer;
        font: inherit;
        line-height: 1;
        padding: 0.4rem 0.6rem;
      }

      button:hover:not(:disabled) {
        border-color: var(--color-marker-background, #006eff);
      }

      button:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      button:disabled {
        cursor: default;
        opacity: 0.5;
      }

      output {
        font-variant-numeric: tabular-nums;
        min-width: 4ch;
        text-align: center;
      }

      .content {
        font-size: var(--resizable-text-size, 100%);
      }
    `

    const toolbar = document.createElement('div')
    toolbar.className = 'toolbar'

    this.decreaseButton = this.createButton('Decrease font size', 'A−', -resizableTextStep)
    this.increaseButton = this.createButton('Increase font size', 'A+', resizableTextStep)
    this.output = document.createElement('output')
    this.output.setAttribute('aria-live', 'polite')
    this.output.setAttribute('aria-label', 'Current font size')
    toolbar.append(this.decreaseButton, this.output, this.increaseButton)

    this.content = document.createElement('div')
    this.content.className = 'content'
    this.content.append(document.createElement('slot'))
    this.shadowRoot.replaceChildren(style, toolbar, this.content)
    this.shadowRoot.addEventListener('click', this.onClick)

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'group')
    }
    if (!this.hasAttribute('aria-label')) {
      this.setAttribute('aria-label', 'Resizable text controls')
    }

    this.loadSize()
    this.renderSize()
  }

  disconnectedCallback() {
    this.shadowRoot.removeEventListener('click', this.onClick)
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name !== 'data-key' || oldValue === newValue || !this.isConnected) {
      return
    }

    this.loadSize()
    this.renderSize()
  }

  createButton(label, text, change) {
    const button = document.createElement('button')
    button.type = 'button'
    button.textContent = text
    button.setAttribute('aria-label', label)
    button.dataset.change = String(change)
    return button
  }

  onClick(event) {
    const button = event.target.closest('button[data-change]')
    if (!button || button.disabled) {
      return
    }

    this.size = Math.min(
      resizableTextMaximum,
      Math.max(resizableTextMinimum, this.size + Number(button.dataset.change)),
    )
    this.renderSize()
    this.saveSize()
  }

  loadSize() {
    const key = this.getAttribute('data-key')
    if (!key) {
      this.size = resizableTextDefault
      return
    }

    try {
      const storedSize = Number(window.localStorage.getItem(key))
      this.size = Number.isFinite(storedSize) && storedSize >= resizableTextMinimum
        ? Math.min(storedSize, resizableTextMaximum)
        : resizableTextDefault
    } catch {
      this.size = resizableTextDefault
    }
  }

  saveSize() {
    const key = this.getAttribute('data-key')
    if (!key) {
      return
    }

    try {
      window.localStorage.setItem(key, String(this.size))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

  renderSize() {
    if (!this.content) {
      return
    }

    this.content.style.setProperty('--resizable-text-size', `${this.size}%`)
    this.output.textContent = `${this.size}%`
    this.decreaseButton.disabled = this.size <= resizableTextMinimum
    this.increaseButton.disabled = this.size >= resizableTextMaximum
  }
}

if (!customElements.get('resizable-text')) {
  customElements.define('resizable-text', ResizableText)
}
