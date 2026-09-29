// use as <timer-start-pause><presentation-timer></presentation-timer></timer-start-pause>

class TimerStartPause extends HTMLElement {
  constructor() {
    super()
    // critical! otherwise "this" will refer to event context
    // and the timerElement will be <null>
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onTimerStateChange = this.onTimerStateChange.bind(this)
    this.onTimerTick = this.onTimerTick.bind(this)
    this.timerElement = null
    this.timerStateObserver = new MutationObserver(this.onTimerStateChange)
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
    window.removeEventListener(navigationCompleteEventType, this.onNavigationComplete)
    this.timerElement?.removeEventListener('presentation-timer-tick', this.onTimerTick)
    this.timerStateObserver.disconnect()
    this.timerElement = null
  }

  onDOMReady() {
    this.initialize()
  }

  initialize() {
    this.timerElement = this.querySelector('presentation-timer')
    if (!this.timerElement) {
      console.error('timer-start-pause requires a presentation-timer child')
      return
    }

    this.timerElement.addEventListener('presentation-timer-tick', this.onTimerTick)
    this.restoreTimerState()
    this.timerStateObserver.observe(this.timerElement, {
      attributes: true,
      attributeFilter: ['paused'],
    })
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
  }

  onNavigationComplete(event) {
    if (!this.timerElement) {
      return
    }

    const slideIndex = event.detail?.slideIndex
    if (slideIndex === 0) {
      this.timerElement.pause()
    } else if (this.timerElement.hasAttribute('paused')) {
      this.timerElement.start()
    }
  }

  onTimerTick(event) {
    if (event.currentTarget !== this.timerElement) {
      return
    }

    const duration = Number(event.detail?.totalDuration)
    const remainingDuration = Number(event.detail?.remainingDuration)
    if (
      !Number.isFinite(duration) || duration <= 0 ||
      !Number.isFinite(remainingDuration) || remainingDuration < 0
    ) {
      return
    }

    this.storeTimerState({
      duration,
      remainingDuration: Math.min(remainingDuration, duration),
      running: remainingDuration > 0 && this.timerElement.running,
    })
  }

  onTimerStateChange() {
    if (!this.timerElement) {
      return
    }

    const duration = this.timerElement.duration
    const elapsed = this.timerElement.elapsed + (
      this.timerElement.running
        ? performance.now() - this.timerElement.startedAt
        : 0
    )
    this.storeTimerState({
      duration,
      remainingDuration: Math.max(0, duration - elapsed),
      running: this.timerElement.running && elapsed < duration,
    })
  }

  getStorageKey() {
    const presentationID = document.querySelector('html')?.dataset.id
    if (!presentationID || !this.timerElement) {
      return null
    }

    return `${presentationID}:timer:${this.timerElement.id || 'default'}`
  }

  restoreTimerState() {
    const key = this.getStorageKey()
    if (!key) {
      return
    }

    let state
    try {
      const storedState = window.localStorage.getItem(key)
      if (!storedState) {
        return
      }
      state = JSON.parse(storedState)
    } catch {
      return
    }

    const duration = Number(state?.duration)
    const remainingDuration = Number(state?.remainingDuration)
    if (
      !Number.isFinite(duration) || duration <= 0 ||
      !Number.isFinite(remainingDuration) || remainingDuration < 0 ||
      remainingDuration > duration || typeof state.running !== 'boolean'
    ) {
      return
    }

    this.timerElement.pause()
    this.timerElement.duration = duration
    this.timerElement.elapsed = duration - remainingDuration
    this.timerElement.updateProgress(this.timerElement.elapsed)
    this.timerElement.updateState()
    if (state.running && remainingDuration > 0) {
      this.timerElement.start()
    }
  }

  storeTimerState(state) {
    const key = this.getStorageKey()
    if (!key) {
      return
    }

    try {
      window.localStorage.setItem(key, JSON.stringify(state))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }
}

if (!customElements.get('timer-start-pause')) {
  customElements.define('timer-start-pause', TimerStartPause)
}
