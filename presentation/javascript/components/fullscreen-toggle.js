// use as custom element <fullscreen-toggle> with its label as text content

class FullscreenToggle extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onKeyDown = this.onKeyDown.bind(this)
    this.onActivateKeyDown = this.onActivateKeyDown.bind(this)
    this.updateState = this.updateState.bind(this)
    this.onClick = this.onClick.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        cursor: pointer;
        display: inline-block;
      }

      :host([fullscreen]) {
        color: red;
      }

      :host(:focus-visible) {
        outline: 2px solid currentColor;
        outline-offset: 2px;
      }
    `

    const slot = document.createElement('slot')
    this.shadowRoot.replaceChildren(style, slot)

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'button')
    }
    if (!this.hasAttribute('tabindex')) {
      this.setAttribute('tabindex', '0')
    }

    this.addEventListener('click', this.onClick)
    this.addEventListener('keydown', this.onActivateKeyDown)
    document.addEventListener('keydown', this.onKeyDown)
    document.addEventListener('fullscreenchange', this.updateState)
    document.addEventListener('webkitfullscreenchange', this.updateState)
    this.updateState()
  }

  disconnectedCallback() {
    this.removeEventListener('click', this.onClick)
    this.removeEventListener('keydown', this.onActivateKeyDown)
    document.removeEventListener('keydown', this.onKeyDown)
    document.removeEventListener('fullscreenchange', this.updateState)
    document.removeEventListener('webkitfullscreenchange', this.updateState)
  }

  onClick() {
    this.toggleFullscreen()
  }

  onActivateKeyDown(event) {
    if (
      (event.code !== 'Enter' && event.code !== 'Space') ||
      this.isEditableTarget(event.target)
    ) {
      return
    }

    event.preventDefault()
    this.toggleFullscreen()
  }

  onKeyDown(event) {
    if (
      event.code !== 'KeyF' ||
      event.altKey ||
      event.ctrlKey ||
      event.metaKey ||
      event.shiftKey ||
      this.isEditableTarget(event.target)
    ) {
      return
    }

    event.preventDefault()
    this.toggleFullscreen()
  }

  isEditableTarget(target) {
    return target instanceof Element && (
      target.isContentEditable ||
      target.matches('input, textarea, select')
    )
  }

  isFullscreen() {
    return Boolean(
      document.fullscreenElement ||
      document.webkitFullscreenElement ||
      document.fullScreenElement ||
      document.webkitIsFullScreen
    )
  }

  async toggleFullscreen() {
    try {
      if (this.isFullscreen()) {
        const exit =
          document.exitFullscreen ||
          document.exitFullScreen ||
          document.cancelFullScreen ||
          document.webkitCancelFullScreen
        if (exit) {
          await exit.call(document)
        }
        return
      }

      const root = document.documentElement
      const request =
        root.requestFullscreen ||
        root.requestFullScreen ||
        root.webkitRequestFullScreen
      if (request) {
        await request.call(root)
      }
    } catch (error) {
      this.dispatchEvent(new CustomEvent('fullscreenerror', {
        bubbles: true,
        composed: true,
        detail: error,
      }))
    }
  }

  updateState() {
    const fullscreen = this.isFullscreen()
    this.toggleAttribute('fullscreen', fullscreen)
    this.setAttribute('aria-pressed', String(fullscreen))
  }
}

if (!customElements.get('fullscreen-toggle')) {
  customElements.define('fullscreen-toggle', FullscreenToggle)
}
