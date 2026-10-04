// ============== synchronize.js =====================

class WindowState {
  DocumentID
  WindowID
  SlideIndex
  ConcealedListItemCount
  TimerElapsedDuration
  TimerUtmostDuration
  IsTimerRunning
  IsCurtainDown
  IsReloading
}

class Synchronizer {
  DocumentID
  WindowID
  Curtain
  MainTimer
  Channel

  constructor() {
    this.DocumentID = document.querySelector("html").dataset.id ||
      "random"+Math.random().toString(36) + Date.now().toString(36)
    this.WindowID = Math.random().toString(36) + '|' + Date.now()
    this.Curtain = document.querySelector('presentation-curtain')
    this.MainTimer = document.querySelector('presentation-timer#mainTimer')
    this.Channel = new BroadcastChannel(this.DocumentID)

    this.broadcastWindowState = debounce(this.broadcastWindowState.bind(this), 120)
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onCurtainToggle = this.onCurtainToggle.bind(this)
    this.onFocus = this.onFocus.bind(this)
    this.onBlur = this.onBlur.bind(this)
    this.onMessage = this.onMessage.bind(this)

    this.setWindowFocused(document.hasFocus())
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    window.addEventListener(curtainToggleEventType, this.onCurtainToggle)
    window.addEventListener('focus', this.onFocus)
    window.addEventListener('blur', this.onBlur)
    this.Channel.addEventListener('message', this.onMessage)
    this.broadcastWindowState()
  }


  setMainTimerState(state) {
    const timer = this.MainTimer
    const duration = Number(state?.TimerUtmostDuration)
    const elapsed = Number(state?.TimerElapsedDuration)
    if (
      !timer || state?.DocumentID !== this.DocumentID ||
      !Number.isFinite(duration) || duration <= 0 ||
      !Number.isFinite(elapsed) || elapsed < 0 || elapsed > duration ||
      typeof state.IsTimerRunning !== 'boolean'
    ) {
      return
    }

    const expired = elapsed >= duration
    const running = state.IsTimerRunning && !expired && state.IsCurtainDown !== true
    timer.pause()
    timer.duration = duration
    timer.elapsed = elapsed
    timer.setExpired(expired)
    timer.updateProgress(timer.elapsed)
    timer.updateState()
    if (running) {
      timer.start()
    }
  }

  setMainCurtainState(state) {
    if (this.Curtain && typeof state?.IsCurtainDown === 'boolean') {
      this.Curtain.toggleAttribute('open', state.IsCurtainDown)
    }
  }

  broadcastWindowState(reloading = false) {
    const state = new WindowState()
    state.DocumentID = this.DocumentID
    state.WindowID = this.WindowID
    state.SlideIndex = currentSlide
    state.ConcealedListItemCount = getCurrentConcealedListItems().length
    state.IsCurtainDown = Boolean(this.Curtain?.hasAttribute('open'))
    state.IsReloading = reloading === true

    const timer = this.MainTimer
    const duration = Number(timer?.duration)
    const elapsed = Number(timer?.elapsed) + (
      timer?.running
        ? performance.now() - Number(timer.startedAt)
        : 0
    )
    if (timer && Number.isFinite(duration) && duration > 0 && Number.isFinite(elapsed)) {
      const expired = timer.hasAttribute('expired') || elapsed >= duration
      state.TimerElapsedDuration = expired
        ? duration
        : Math.max(0, Math.min(duration, elapsed))
      state.TimerUtmostDuration = duration
      state.IsTimerRunning = timer.running && !expired
    } else {
      state.TimerElapsedDuration = 0
      state.TimerUtmostDuration = 0
      state.IsTimerRunning = false
    }

    this.Channel.postMessage(state)
  }

  setWindowFocused(focused) {
    document.body.classList.toggle('is-focused', focused)
  }

  onNavigationComplete() {
    this.broadcastWindowState()
  }

  onCurtainToggle(event) {
    if (event.target !== this.Curtain) {
      return
    }
    this.broadcastWindowState()
  }

  onFocus() {
    this.setWindowFocused(true)
    this.broadcastWindowState()
  }

  onBlur() {
    this.setWindowFocused(false)
    this.broadcastWindowState()
  }

  onMessage(event) {
    const broadcast = event.data
    this.setMainCurtainState(broadcast)
    this.setMainTimerState(broadcast)
    this.setWindowFocused(false)
    if (!broadcast || broadcast.DocumentID !== this.DocumentID) {
      return
    }
    if (broadcast.IsReloading === true) {
      window.location.reload()
      return
    }
    if (broadcast.WindowID === this.WindowID) {
      return
    }

    if (broadcast.SlideIndex !== currentSlide) {
      navigate(broadcast.SlideIndex)
    }
    const concealedListItems = getCurrentConcealedListItems()
    let revealCount = concealedListItems.length - broadcast.ConcealedListItemCount
    for (const element of concealedListItems) {
      if (revealCount === 0) {
        return
      }
      element.classList.add("is-revealed")
      revealCount--
    }
  }
}

const synchronizer = new Synchronizer()
