// use as <timer-start-pause><presentation-timer></presentation-timer></timer-start-pause>

const mainTimerLocalStorageKey = "mainTimerLocalStorageKey"

class TimerState {
  ElapsedDuration
  UtmostDuration
  IsRunning
}

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
      attributeFilter: ['paused'],
    })
    this.#timer.addEventListener('change', this.onTimerChange)
    this.addEventListener('presentation-timer-tick', this.onTimerStateChange)
    this.addEventListener('click', this.onTimerStateChange)
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.addEventListener('change', this.onCurtainChange)

    const currentState = this.getTimerState()
    const storedState = this.readStoredTimerState()
    const initialState = storedState ?? currentState
    if (!storedState) {
      initialState.IsRunning = false
    }
    if (this.#curtain) {
      initialState.IsRunning = !this.#curtain.IsDown()
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
      IsRunning: slideIndex !== 0 && !this.isCurtainDown(),
    })
  }

  onCurtainChange(event) {
    if (event.target !== this.#curtain) {
      this.#curtain = document.querySelector('presentation-curtain')
      if (event.target !== this.#curtain) {
        return
      }
    }

    this.setTimerState({ IsRunning: !this.#curtain.IsDown() })
  }

  isCurtainDown() {
    return Boolean(this.#curtain?.IsDown())
  }

  onTimerStateChange() {
    this.setTimerState()
  }

  onTimerChange(event) {
    this.onTimerStateChange()
  }

  setTimerState(state) {
    if (state) {
      this.pendingTimerState = {
        ...this.pendingTimerState,
        ...state,
      }
    }

    window.clearTimeout(this.timerStateTimeout)
    this.timerStateTimeout = window.setTimeout(() => {
      this.timerStateTimeout = null
      const requestedState = this.pendingTimerState
      this.pendingTimerState = null

      if (!requestedState) {
        this.storeTimerState(this.getTimerState())
        return
      }

      const targetState = Object.assign(
        new TimerState(),
        this.getTimerState(),
        requestedState,
      )
      const elapsed = Number(targetState.ElapsedDuration)
      const utmost = Number(targetState.UtmostDuration)
      if (
        !Number.isFinite(elapsed) || elapsed < 0 ||
        !Number.isFinite(utmost) || utmost <= 0 ||
        typeof targetState.IsRunning !== 'boolean'
      ) {
        return
      }

      const timer = this.#timer
      this.timerStateObserver.disconnect()
      try {
        timer.setDuration(elapsed, utmost)
        if (targetState.IsRunning) {
          timer.start()
        }
      } finally {
        if (this.isConnected) {
          this.timerStateObserver.observe(this, {
            attributes: true,
            subtree: true,
            attributeFilter: ['paused'],
          })
        }
      }

      this.storeTimerState(this.getTimerState())
    }, 600)
  }

  getTimerState() {
    const state = new TimerState()
    state.ElapsedDuration = Number(this.#timer.ElapsedDuration)
    state.UtmostDuration = Number(this.#timer.UtmostDuration)
    state.IsRunning = this.#timer.TimerRunning
    return state
  }

  readStoredTimerState() {
    let state
    try {
      const storedState = window.localStorage.getItem(mainTimerLocalStorageKey)
      if (!storedState) {
        return null
      }
      state = JSON.parse(storedState)
    } catch {
      return null
    }

    const elapsed = Number(state?.ElapsedDuration)
    const utmost = Number(state?.UtmostDuration)
    if (
      !Number.isFinite(elapsed) || elapsed < 0 ||
      !Number.isFinite(utmost) || utmost <= 0 ||
      typeof state?.IsRunning !== 'boolean'
    ) {
      return null
    }

    const timerState = new TimerState()
    timerState.ElapsedDuration = elapsed
    timerState.UtmostDuration = utmost
    timerState.IsRunning = state.IsRunning
    return timerState
  }

  storeTimerState(state) {
    try {
      window.localStorage.setItem(mainTimerLocalStorageKey, JSON.stringify(
        Object.assign(new TimerState(), state),
      ))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

}

if (!customElements.get('timer-start-pause')) {
  customElements.define('timer-start-pause', TimerStartPause)
}
