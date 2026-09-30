// use as <timer-start-pause><presentation-timer></presentation-timer></timer-start-pause>

class TimerStartPause extends HTMLElement {
  constructor() {
    super()
    // critical! otherwise "this" will refer to event context
    // and the timerElement will be <null>
    this.onDOMReady = this.onDOMReady.bind(this)
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onCurtainToggle = this.onCurtainToggle.bind(this)
    this.onTimerStateChange = this.onTimerStateChange.bind(this)
    this.onTimerTick = this.onTimerTick.bind(this)
    this.onTimerStateBroadcast = this.onTimerStateBroadcast.bind(this)
    this.sourceID = `${Math.random().toString(36).slice(2)}-${Date.now().toString(36)}`
    this.timerElement = null
    this.timerStateChannel = null
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
    window.removeEventListener(curtainToggleEventType, this.onCurtainToggle)
    this.timerElement?.removeEventListener('presentation-timer-tick', this.onTimerTick)
    this.timerStateChannel?.removeEventListener('message', this.onTimerStateBroadcast)
    this.timerStateChannel?.close()
    this.timerStateChannel = null
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
    this.timerStateChannel = new BroadcastChannel(`${documentID}TimerState`)
    this.timerStateChannel.addEventListener('message', this.onTimerStateBroadcast)
    this.restoreTimerState()
    this.timerStateObserver.observe(this.timerElement, {
      attributes: true,
      attributeFilter: ['paused', 'expired'],
    })
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.addEventListener(curtainToggleEventType, this.onCurtainToggle)
    this.syncTimerWithCurtain()
    this.onTimerStateChange()
  }

  onNavigationComplete(event) {
    if (!this.timerElement) {
      return
    }

    const slideIndex = event.detail?.slideIndex
    this.setTimerRunning(slideIndex !== 0 && !this.isCurtainOpen())
  }

  onCurtainToggle(event) {
    if (
      !this.timerElement ||
      event.target !== document.querySelector('presentation-curtain')
    ) {
      return
    }

    this.setTimerRunning(event.detail?.open !== true)
  }

  syncTimerWithCurtain() {
    const curtain = document.querySelector('presentation-curtain')
    if (curtain) {
      this.setTimerRunning(!curtain.hasAttribute('open'))
    }
  }

  isCurtainOpen() {
    const curtain = document.querySelector('presentation-curtain')
    return Boolean(curtain?.hasAttribute('open'))
  }

  setTimerRunning(shouldRun) {
    if (!this.timerElement || this.timerElement.hasAttribute('expired')) {
      return
    }

    if (shouldRun && this.timerElement.hasAttribute('paused')) {
      this.timerElement.start()
    } else if (!shouldRun && !this.timerElement.hasAttribute('paused')) {
      this.timerElement.pause()
    }
    this.onTimerStateChange()
  }

  onTimerStateBroadcast(event) {
    const broadcast = event.data
    if (
      !this.timerElement ||
      !broadcast ||
      typeof broadcast.sourceID !== 'string' ||
      broadcast.sourceID === this.sourceID ||
      broadcast.timerID !== (this.timerElement.id || 'default') ||
      this.timerElement.hasAttribute('expired')
    ) {
      return
    }

    const state = broadcast.state
    const duration = Number(state?.duration)
    const remainingDuration = Number(state?.remainingDuration)
    if (
      !Number.isFinite(duration) || duration <= 0 ||
      !Number.isFinite(remainingDuration) || remainingDuration < 0 ||
      remainingDuration > duration || typeof state?.running !== 'boolean' ||
      typeof state.expired !== 'boolean'
    ) {
      return
    }

    this.applyTimerState({
      duration,
      remainingDuration,
      running: state.running,
      expired: state.expired,
    })
  }

  applyTimerState(state) {
    if (!this.timerElement || this.timerElement.hasAttribute('expired')) {
      return
    }

    this.timerStateObserver.disconnect()
    try {
      this.timerElement.pause()
      this.timerElement.duration = state.duration
      this.timerElement.elapsed = state.duration - state.remainingDuration
      this.timerElement.updateProgress(this.timerElement.elapsed)
      if (state.expired) {
        this.timerElement.setExpired(true)
      }
      this.timerElement.updateState()
      if (state.running && state.remainingDuration > 0 && !state.expired) {
        this.timerElement.start()
      }
    } finally {
      if (this.timerElement) {
        this.timerStateObserver.observe(this.timerElement, {
          attributes: true,
          attributeFilter: ['paused', 'expired'],
        })
      }
    }
  }

  onTimerTick(event) {
    if (
      !this.timerElement ||
      event.currentTarget !== this.timerElement ||
      this.timerElement.hasAttribute('expired')
    ) {
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

    if (this.timerElement.hasAttribute('expired')) {
      this.broadcastTimerState({
        duration: this.timerElement.duration,
        remainingDuration: 0,
        running: false,
        expired: true,
      })
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

  broadcastTimerState(state) {
    try {
      this.timerStateChannel?.postMessage({
        sourceID: this.sourceID,
        timerID: this.timerElement?.id || 'default',
        state,
      })
    } catch {
      // Broadcast channels may be unavailable in restricted browsing contexts.
    }
  }

  getStorageKey() {
    const presentationID = document.querySelector('html')?.dataset.id
    if (!presentationID || !this.timerElement) {
      return null
    }

    return `${presentationID}:timer:${this.timerElement.id || 'default'}`
  }

  restoreTimerState() {
    if (!this.timerElement || this.timerElement.hasAttribute('expired')) {
      return
    }

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
    if (!this.timerElement || this.timerElement.hasAttribute('expired')) {
      return
    }

    const key = this.getStorageKey()
    if (key) {
      try {
        window.localStorage.setItem(key, JSON.stringify(state))
      } catch {
        // Local storage may be unavailable in restricted browsing contexts.
      }
    }

    this.broadcastTimerState({
      ...state,
      expired: false,
    })
  }
}

if (!customElements.get('timer-start-pause')) {
  customElements.define('timer-start-pause', TimerStartPause)
}
