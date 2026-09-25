// use as custom element <fullscreen-toggle> anywhere on the page

class FullscreenToggle extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onKeyDown = this.onKeyDown.bind(this)
    this.updateButton = this.updateButton.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        display: inline-block;
      }

      button {
        align-items: center;
        display: inline-flex;
      }
    `

    this.button = document.createElement('button')
    this.button.type = 'button'
    this.button.addEventListener('click', () => this.toggleFullscreen())
    this.shadowRoot.replaceChildren(style, this.button)

    document.addEventListener('keydown', this.onKeyDown)
    document.addEventListener('fullscreenchange', this.updateButton)
    document.addEventListener('webkitfullscreenchange', this.updateButton)
    this.updateButton()
  }

  disconnectedCallback() {
    document.removeEventListener('keydown', this.onKeyDown)
    document.removeEventListener('fullscreenchange', this.updateButton)
    document.removeEventListener('webkitfullscreenchange', this.updateButton)
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

  updateButton() {
    if (!this.button) {
      return
    }

    const fullscreen = this.isFullscreen()
    this.button.textContent = fullscreen ? 'Exit fullscreen' : 'Fullscreen'
    this.button.setAttribute('aria-label', this.button.textContent)
    this.button.setAttribute('aria-pressed', String(fullscreen))
  }
}

if (!customElements.get('fullscreen-toggle')) {
  customElements.define('fullscreen-toggle', FullscreenToggle)
}
