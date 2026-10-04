// use as <timer-start-pause><presentation-timer></presentation-timer></timer-start-pause>

class TimerStartPause extends HTMLElement {
  #timer
  #curtain

  constructor() {
    super()
    this.#timer = this.querySelector('presentation-timer')
    this.#curtain = document.querySelector('presentation-curtain')
    this.timerStateTimeout = null
    this.pendingTimerState = null
    this.isInitialized = false
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onDOMReady = this.onDOMReady.bind(this)
    this.onCurtainChange = this.onCurtainChange.bind(this)
    this.onTimerStateChange = this.onTimerStateChange.bind(this)
    this.onTimerChange = this.onTimerChange.bind(this)
    this.timerStateObserver = new MutationObserver(this.onTimerStateChange)
  }

  connectedCallback() {
    this.initialize()
  }

  disconnectedCallback() {
    document.removeEventListener('DOMContentLoaded', this.onDOMReady)
    window.removeEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.removeEventListener('change', this.onCurtainChange)
    this.#timer?.removeEventListener('change', this.onTimerChange)
    this.removeEventListener('presentation-timer-tick', this.onTimerStateChange)
    this.removeEventListener('click', this.onTimerStateChange)
    this.timerStateObserver.disconnect()
    window.clearTimeout(this.timerStateTimeout)
    this.timerStateTimeout = null
    this.pendingTimerState = null
    this.isInitialized = false
  }

  initialize() {
    if (this.isInitialized) {
      return
    }
    this.#timer ??= this.querySelector('presentation-timer')
    if (!this.#timer) {
      console.error('timer-start-pause requires a presentation-timer child')
      return
    }

    this.#curtain = document.querySelector('presentation-curtain')
    this.isInitialized = true
    this.timerStateObserver.observe(this, {
      attributes: true,
      subtree: true,
      attributeFilter: ['paused', 'expired'],
    })
    this.#timer.addEventListener('change', this.onTimerChange)
    this.addEventListener('presentation-timer-tick', this.onTimerStateChange)
    this.addEventListener('click', this.onTimerStateChange)
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.addEventListener('change', this.onCurtainChange)

    const currentState = this.getTimerState()
    const storedState = this.readStoredTimerState(currentState.timerID)
    const initialState = {
      ...currentState,
      ...(storedState && !currentState.expired ? storedState : {
        // default values
        running: false,
      }),
    }
    if (this.#curtain && !initialState.expired) {
      initialState.running = !this.#curtain.IsDown()
    }
    this.setTimerState(initialState)

    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', this.onDOMReady, { once: true })
    } else {
      this.onDOMReady()
    }
  }

  onDOMReady() {
    this.#curtain = document.querySelector('presentation-curtain')
    if (typeof currentSlide !== 'number') {
      return
    }
    this.onNavigationComplete({ detail: { slideIndex: currentSlide } })
  }

  onNavigationComplete(event) {
    const slideIndex = event.detail?.slideIndex
    this.setTimerState({
      running: slideIndex !== 0 && !this.isCurtainDown(),
    })
  }

  onCurtainChange(event) {
    if (event.target !== this.#curtain) {
      this.#curtain = document.querySelector('presentation-curtain')
      if (event.target !== this.#curtain) {
        return
      }
    }

    this.setTimerState({ running: !this.#curtain.IsDown() })
  }

  isCurtainDown() {
    return Boolean(this.#curtain?.IsDown())
  }

  onTimerStateChange() {
    this.setTimerState()
  }

  onTimerChange(event) {
    if (event.target !== this.#timer) {
      return
    }

    this.querySelector('timer-set')?.updateDurationFromTimer?.(this.#timer)
    this.onTimerStateChange()
  }

  setTimerState(state) {
    if (state) {
      const nextState = { ...state }
      if (
        Number.isFinite(nextState.remainingDuration) &&
        !Number.isFinite(nextState.capturedAt)
      ) {
        nextState.capturedAt = performance.now()
        nextState.counting = nextState.running === true
      }
      this.pendingTimerState = {
        ...this.pendingTimerState,
        ...nextState,
      }
    }

    window.clearTimeout(this.timerStateTimeout)
    this.timerStateTimeout = window.setTimeout(() => {
      this.timerStateTimeout = null
      const requestedState = this.pendingTimerState
      this.pendingTimerState = null

      if (!requestedState) {
        const currentState = this.getTimerState()
        this.storeTimerState(currentState)
        return
      }

      const currentState = this.getTimerState()
      const targetState = { ...currentState, ...requestedState }
      const duration = Number(targetState.duration)
      let remainingDuration = Number(targetState.remainingDuration)
      if (targetState.counting === true && Number.isFinite(targetState.capturedAt)) {
        remainingDuration -= performance.now() - targetState.capturedAt
      }
      if (
        !Number.isFinite(duration) || duration <= 0 ||
        !Number.isFinite(remainingDuration)
      ) {
        return
      }

      if (targetState.expired === true) {
        remainingDuration = 0
      }
      remainingDuration = Math.min(duration, Math.max(0, remainingDuration))
      const expired = targetState.expired === true || remainingDuration === 0
      const running = targetState.running === true && remainingDuration > 0 && !expired
      const timer = this.#timer

      this.timerStateObserver.disconnect()
      try {
        timer.setDuration(duration - remainingDuration, duration)
        if (running) {
          timer.start()
        }
      } finally {
        if (this.isConnected) {
          this.timerStateObserver.observe(this, {
            attributes: true,
            subtree: true,
            attributeFilter: ['paused', 'expired'],
          })
        }
      }

      const updatedState = this.getTimerState()
      this.storeTimerState(updatedState)
    }, 600)
  }

  getTimerState() {
    const timer = this.#timer
    const now = performance.now()
    const duration = Number(timer.UtmostDuration)
    const elapsed = Number(timer.ElapsedDuration)
    const expired = timer.hasAttribute('expired') || elapsed >= duration

    return {
      timerID: timer.id || 'default',
      duration,
      remainingDuration: expired ? 0 : Math.max(0, duration - elapsed),
      running: timer.TimerRunning && !expired,
      expired,
      capturedAt: now,
      counting: timer.TimerRunning && !expired,
    }
  }

  getStorageKey(timerID) {
    const presentationID = document.querySelector('html')?.dataset.id
    if (!presentationID || !timerID) {
      return null
    }

    return `${presentationID}:timer:${timerID}`
  }

  readStoredTimerState(timerID) {
    const key = this.getStorageKey(timerID)
    if (!key) {
      return null
    }

    let state
    try {
      const storedState = window.localStorage.getItem(key)
      if (!storedState) {
        return null
      }
      state = JSON.parse(storedState)
    } catch {
      return null
    }

    const duration = Number(state?.duration)
    const remainingDuration = Number(state?.remainingDuration)
    if (
      !Number.isFinite(duration) || duration <= 0 ||
      !Number.isFinite(remainingDuration) || remainingDuration < 0 ||
      remainingDuration > duration || typeof state.running !== 'boolean' ||
      (state.expired !== undefined && typeof state.expired !== 'boolean')
    ) {
      return null
    }

    return {
      duration,
      remainingDuration,
      running: state.running,
      expired: state.expired === true,
    }
  }

  storeTimerState(state) {
    const key = this.getStorageKey(state.timerID)
    if (!key) {
      return
    }

    try {
      window.localStorage.setItem(key, JSON.stringify({
        duration: state.duration,
        remainingDuration: state.remainingDuration,
        running: state.running,
        expired: state.expired,
      }))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

}

if (!customElements.get('timer-start-pause')) {
  customElements.define('timer-start-pause', TimerStartPause)
}
