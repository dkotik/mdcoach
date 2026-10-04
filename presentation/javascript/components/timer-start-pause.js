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
    this.isInitialized = false
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onDOMReady = this.onDOMReady.bind(this)
    this.onCurtainChange = this.onCurtainChange.bind(this)
    this.onTimerStateChange = this.onTimerStateChange.bind(this)
  }

  connectedCallback() {
    this.initialize()
  }

  disconnectedCallback() {
    document.removeEventListener('DOMContentLoaded', this.onDOMReady)
    window.removeEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.removeEventListener('change', this.onCurtainChange)
    this.#timer?.removeEventListener('change', this.onTimerStateChange)
    this.removeEventListener('presentation-timer-tick', this.onTimerStateChange)
    this.removeEventListener('click', this.onTimerStateChange)
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
    this.#timer.addEventListener('change', this.onTimerStateChange)
    this.addEventListener('presentation-timer-tick', this.onTimerStateChange)
    this.addEventListener('click', this.onTimerStateChange)
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.addEventListener('change', this.onCurtainChange)

    const initialState = new TimerState()
    try {
      const storedState = window.localStorage.getItem(mainTimerLocalStorageKey)
      if (storedState) {
        Object.assign(initialState, JSON.parse(storedState))
      }
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }

    const elapsed = Number(initialState.ElapsedDuration)
    const utmost = Number(initialState.UtmostDuration)
    if (
      !Number.isFinite(elapsed) || elapsed < 0 ||
      !Number.isFinite(utmost) || utmost <= 0 ||
      typeof initialState.IsRunning !== 'boolean'
    ) {
      initialState.ElapsedDuration = Number(this.#timer.ElapsedDuration) || 0
      initialState.UtmostDuration = Number(this.#timer.UtmostDuration) || this.#timer.getDuration()
      initialState.IsRunning = false
    } else {
      initialState.ElapsedDuration = elapsed
      initialState.UtmostDuration = utmost
    }
    if (this.#curtain) {
      initialState.IsRunning = !this.#curtain.IsDown()
    }
    this.#timer.setDuration(initialState.ElapsedDuration, initialState.UtmostDuration)
    if (initialState.IsRunning) {
      this.#timer.start()
    }
    this.onTimerStateChange()

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
    if (slideIndex === 0) {
      this.#timer.pause()
    } else if (!this.#timer.IsRunning && !this.#curtain?.IsDown()) {
      this.#timer.start()
    }
    this.onTimerStateChange()
  }

  onCurtainChange(event) {
    if (event.target !== this.#curtain) {
      return
    }

    if (this.#curtain.IsDown()) {
      this.#timer.pause()
    } else {
      this.#timer.start()
    }
    this.onTimerStateChange()
  }

  onTimerStateChange() {
    const state = new TimerState()
    state.ElapsedDuration = Number(this.#timer.ElapsedDuration)
    state.UtmostDuration = Number(this.#timer.UtmostDuration)
    state.IsRunning = this.#timer.IsRunning

    // console.log("timerState:", state)

    try {
      window.localStorage.setItem(mainTimerLocalStorageKey, JSON.stringify(state))
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

}

if (!customElements.get('timer-start-pause')) {
  customElements.define('timer-start-pause', TimerStartPause)
}
