// use as <timer-start-pause><presentation-timer></presentation-timer></timer-start-pause>

class TimerStartPause extends HTMLElement {
  constructor() {
    super()
    // critical! otherwise "this" will refer to event context
    // and the timerElement will be <null>
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.timerElement = null
  }

  connectedCallback() {
    this.timerElement = this.querySelector('presentation-timer')
    if (!this.timerElement) {
      console.error("there is no child timer")
    }
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
  }

  disconnectedCallback() {
    window.removeEventListener(navigationCompleteEventType, this.onNavigationComplete)
    this.timerElement = null
  }

  onNavigationComplete(event) {
    const slideNumber = event.detail.slideIndex
    if (slideNumber === 0) {
      this.timerElement.pause()
    } else if (this.timerElement.hasAttribute('paused')) {
      this.timerElement.start()
    }
  }
}

if (!customElements.get('timer-start-pause')) {
  customElements.define('timer-start-pause', TimerStartPause)
}
