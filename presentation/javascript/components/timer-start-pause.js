// use as <timer-start-pause><presentation-timer></presentation-timer></timer-start-pause>

class TimerStartPause extends HTMLElement {
  constructor() {
    super()
    this.timerElement = this.querySelector('presentation-timer')
    this.hasTimerElement = Boolean(this.timerElement)
    this.timerStateTimeout = null
    this.pendingTimerState = null
    this.isInitialized = false
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onCurtainToggle = this.onCurtainToggle.bind(this)
    this.onTimerStateChange = this.onTimerStateChange.bind(this)
    this.timerStateObserver = new MutationObserver(this.onTimerStateChange)
  }

  connectedCallback() {
    this.initialize()
  }

  disconnectedCallback() {
    window.removeEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.removeEventListener(curtainToggleEventType, this.onCurtainToggle)
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
    if (!this.hasTimerElement) {
      console.error('timer-start-pause requires a presentation-timer child')
      return
    }

    this.isInitialized = true
    this.timerStateObserver.observe(this, {
      attributes: true,
      subtree: true,
      attributeFilter: ['paused', 'expired'],
    })
    this.addEventListener('presentation-timer-tick', this.onTimerStateChange)
    this.addEventListener('click', this.onTimerStateChange)
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.addEventListener(curtainToggleEventType, this.onCurtainToggle)

    const currentState = this.getTimerState()
    const storedState = this.readStoredTimerState(currentState.timerID)
    const initialState = {
      ...currentState,
      ...(storedState && !currentState.expired ? storedState : {}),
    }
    const curtain = document.querySelector('presentation-curtain')
    if (curtain && !initialState.expired) {
      initialState.running = !curtain.hasAttribute('open')
    }
    this.setTimerState(initialState)
  }

  onNavigationComplete(event) {
    const slideIndex = event.detail?.slideIndex
    this.setTimerState({
      running: slideIndex !== 0 && !this.isCurtainOpen(),
    })
  }

  onCurtainToggle(event) {
    if (event.target !== document.querySelector('presentation-curtain')) {
      return
    }

    this.setTimerState({ running: event.detail?.open !== true })
  }

  isCurtainOpen() {
    const curtain = document.querySelector('presentation-curtain')
    return Boolean(curtain?.hasAttribute('open'))
  }

  onTimerStateChange() {
    this.setTimerState()
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
      const timer = this.timerElement

      this.timerStateObserver.disconnect()
      try {
        timer.pause()
        timer.duration = duration
        timer.elapsed = duration - remainingDuration
        timer.setExpired(expired)
        timer.updateProgress(timer.elapsed)
        timer.updateState()
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
    const timer = this.timerElement
    const now = performance.now()
    const duration = Number(timer.duration)
    const elapsed = Number(timer.elapsed) + (
      timer.running
        ? now - Number(timer.startedAt)
        : 0
    )
    const expired = timer.hasAttribute('expired') || elapsed >= duration

    return {
      timerID: timer.id || 'default',
      duration,
      remainingDuration: expired ? 0 : Math.max(0, duration - elapsed),
      running: timer.running && !expired,
      expired,
      capturedAt: now,
      counting: timer.running && !expired,
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
