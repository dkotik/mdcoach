// use as <presentation-play> to present slides on a connected external display

const presentationPlayReceiverParameter = 'presentation-receiver'
const presentationPlayMessageType = 'presentation-slide-state'

class PresentationPlay extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onClick = this.onClick.bind(this)
    this.onConnectionAvailable = this.onConnectionAvailable.bind(this)
    this.onConnectionStateChange = this.onConnectionStateChange.bind(this)
    this.onConnectionError = this.onConnectionError.bind(this)
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onReceiverConnectionAvailable = this.onReceiverConnectionAvailable.bind(this)
    this.onReceiverMessage = this.onReceiverMessage.bind(this)
    this.onDOMReady = this.onDOMReady.bind(this)
    this.connection = undefined
    this.isStarting = false
    this.isDOMReady = document.readyState !== 'loading'
    this.pendingState = undefined
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        display: inline-block;
      }

      :host([hidden]) {
        display: none;
      }

      button {
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-body-subtext, #aaa);
        border-radius: 0.35rem;
        color: var(--color-menu-text, #666);
        cursor: pointer;
        font: inherit;
        padding: 0.5rem 0.75rem;
        white-space: nowrap;
      }

      button:hover:not(:disabled) {
        border-color: var(--color-marker-background, #006eff);
      }

      button:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      button:disabled {
        cursor: not-allowed;
        opacity: 0.55;
      }
    `

    this.button = document.createElement('button')
    this.button.type = 'button'
    this.button.textContent = 'Present'
    this.button.setAttribute('aria-label', 'Present on an external display')
    this.button.addEventListener('click', this.onClick)
    this.shadowRoot.replaceChildren(style, this.button)

    const receiverMode = new URL(window.location.href).searchParams.has(
      presentationPlayReceiverParameter,
    )
    if (receiverMode) {
      this.enterReceiverMode()
      return
    }

    this.setupPresentationRequest()
    window.addEventListener('slideNavigationFinished', this.onNavigationComplete)
  }

  disconnectedCallback() {
    this.button?.removeEventListener('click', this.onClick)
    window.removeEventListener('slideNavigationFinished', this.onNavigationComplete)
    this.request?.removeEventListener('connectionavailable', this.onConnectionAvailable)
    document.removeEventListener('DOMContentLoaded', this.onDOMReady)

    if (this.receiverConnectionList) {
      this.receiverConnectionList.removeEventListener(
        'connectionavailable',
        this.onReceiverConnectionAvailable,
      )
    }

    for (const connection of this.receiverConnections || []) {
      connection.removeEventListener('message', this.onReceiverMessage)
    }
  }

  setupPresentationRequest() {
    if (typeof window.PresentationRequest !== 'function') {
      this.button.disabled = true
      this.button.title = 'Presentation API is not supported by this browser.'
      return
    }

    const url = new URL(this.getAttribute('url') || window.location.href, window.location.href)
    url.searchParams.set(presentationPlayReceiverParameter, '1')
    const slides = Array.from(document.querySelectorAll('main > section'))
    const focusedSlideIndex = slides.findIndex((slide) => slide.classList.contains('is-focused'))
    if (focusedSlideIndex >= 0) {
      url.hash = String(focusedSlideIndex + 1)
    }

    try {
      this.request = new window.PresentationRequest([url.href])
      this.request.addEventListener('connectionavailable', this.onConnectionAvailable)
    } catch (error) {
      this.button.disabled = true
      this.button.title = 'Unable to create a presentation request.'
      this.dispatchError(error)
    }
  }

  async onClick() {
    if (this.isStarting) {
      return
    }

    if (this.connection?.state === 'connected') {
      try {
        this.connection.terminate()
      } catch (error) {
        this.dispatchError(error)
      }
      return
    }

    if (!this.request) {
      return
    }

    this.isStarting = true
    this.updateButton()
    try {
      const connection = await this.request.start()
      this.attachSenderConnection(connection)
    } catch (error) {
      this.dispatchError(error)
    } finally {
      this.isStarting = false
      this.updateButton()
    }
  }

  onConnectionAvailable(event) {
    this.attachSenderConnection(event.connection)
  }

  attachSenderConnection(connection) {
    if (this.connection === connection) {
      return
    }

    this.connection = connection
    connection.addEventListener('connect', this.onConnectionStateChange)
    connection.addEventListener('close', this.onConnectionStateChange)
    connection.addEventListener('terminate', this.onConnectionStateChange)
    connection.addEventListener('statechange', this.onConnectionStateChange)
    connection.addEventListener('error', this.onConnectionError)
    this.updateButton()
    this.sendCurrentState()
  }

  onConnectionStateChange(event) {
    const connection = event.currentTarget
    if (connection.state === 'closed' || connection.state === 'terminated') {
      if (this.connection === connection) {
        this.connection = undefined
      }
    } else if (connection.state === 'connected') {
      this.sendCurrentState()
    }
    this.updateButton()
  }

  onConnectionError(event) {
    this.dispatchError(event.error || event)
  }

  onNavigationComplete() {
    this.sendCurrentState()
  }

  sendCurrentState() {
    if (this.connection?.state !== 'connected') {
      return
    }

    const slides = Array.from(document.querySelectorAll('main > section'))
    const slideIndex = slides.findIndex((slide) => slide.classList.contains('is-focused'))
    if (slideIndex < 0) {
      return
    }

    const listItems = slides[slideIndex].querySelectorAll(
      ':scope > .grid > .content > ul > li',
    )
    const revealedListItemCount = Array.from(listItems)
      .filter((item) => item.classList.contains('is-revealed'))
      .length

    try {
      this.connection.send(JSON.stringify({
        type: presentationPlayMessageType,
        slideIndex,
        revealedListItemCount,
      }))
    } catch (error) {
      this.dispatchError(error)
    }
  }

  enterReceiverMode() {
    this.hidden = true
    document.documentElement.setAttribute('data-presentation-receiver', '')
    this.addReceiverStyles()

    if (!this.isDOMReady) {
      document.addEventListener('DOMContentLoaded', this.onDOMReady, { once: true })
    }

    const receiver = navigator.presentation?.receiver
    if (!receiver) {
      return
    }

    receiver.connectionList.then((connectionList) => {
      this.receiverConnectionList = connectionList
      for (const connection of connectionList.connections) {
        this.attachReceiverConnection(connection)
      }
      connectionList.addEventListener(
        'connectionavailable',
        this.onReceiverConnectionAvailable,
      )
    }).catch((error) => this.dispatchError(error))
  }

  addReceiverStyles() {
    if (document.getElementById('presentation-receiver-styles')) {
      return
    }

    const style = document.createElement('style')
    style.id = 'presentation-receiver-styles'
    style.textContent = `
      html[data-presentation-receiver] presentation-menu,
      html[data-presentation-receiver] presentation-side-button,
      html[data-presentation-receiver] presentation-play,
      html[data-presentation-receiver] presentation-curtain,
      html[data-presentation-receiver] presentation-timer,
      html[data-presentation-receiver] timer-set,
      html[data-presentation-receiver] keystroke-combo,
      html[data-presentation-receiver] slide-notes {
        display: none !important;
      }

      html[data-presentation-receiver] body {
        cursor: none;
      }
    `
    document.head.append(style)
  }

  onReceiverConnectionAvailable(event) {
    this.attachReceiverConnection(event.connection)
  }

  attachReceiverConnection(connection) {
    this.receiverConnections ||= new Set()
    if (this.receiverConnections.has(connection)) {
      return
    }

    this.receiverConnections.add(connection)
    connection.addEventListener('message', this.onReceiverMessage)
  }

  onReceiverMessage(event) {
    if (typeof event.data !== 'string') {
      return
    }

    let state
    try {
      state = JSON.parse(event.data)
    } catch {
      return
    }
    if (state?.type !== presentationPlayMessageType) {
      return
    }

    this.pendingState = state
    if (this.isDOMReady) {
      this.applyState(state)
    }
  }

  onDOMReady() {
    this.isDOMReady = true
    if (this.pendingState) {
      this.applyState(this.pendingState)
    }
  }

  applyState(state) {
    const slides = Array.from(document.querySelectorAll('main > section'))
    if (!Number.isInteger(state.slideIndex) || state.slideIndex < 0 || state.slideIndex >= slides.length) {
      return
    }

    if (typeof navigate === 'function' && typeof currentSlide === 'number') {
      if (currentSlide !== state.slideIndex) {
        navigate(state.slideIndex)
      }
    } else {
      slides.forEach((slide, index) => {
        slide.classList.toggle('is-focused', index === state.slideIndex)
      })
    }

    const listItems = slides[state.slideIndex].querySelectorAll(
      ':scope > .grid > .content > ul > li',
    )
    const revealedListItemCount = Number.isInteger(state.revealedListItemCount)
      ? Math.min(listItems.length, Math.max(0, state.revealedListItemCount))
      : 0
    Array.from(listItems).forEach((item, index) => {
      item.classList.toggle('is-revealed', index < revealedListItemCount)
    })

    window.dispatchEvent(new CustomEvent('slideNavigationFinished', {
      detail: {
        slideIndex: state.slideIndex,
        finalSlideIndex: slides.length - 1,
        concealedListItemCount: listItems.length - revealedListItemCount,
      },
    }))
  }

  updateButton() {
    if (!this.button) {
      return
    }

    const state = this.connection?.state
    const isConnected = state === 'connected'
    const isConnecting = this.isStarting || state === 'connecting'
    this.button.textContent = isConnecting
      ? 'Connecting…'
      : isConnected
        ? 'Stop presenting'
        : 'Present'
    this.button.disabled = isConnecting || !this.request
    this.button.setAttribute('aria-pressed', String(isConnected))
    this.button.setAttribute('aria-label', isConnected
      ? 'Stop presenting on the external display'
      : 'Present on an external display')
  }

  dispatchError(error) {
    this.dispatchEvent(new CustomEvent('presentationerror', {
      bubbles: true,
      composed: true,
      detail: error,
    }))
  }
}

if (!customElements.get('presentation-play')) {
  customElements.define('presentation-play', PresentationPlay)
}
