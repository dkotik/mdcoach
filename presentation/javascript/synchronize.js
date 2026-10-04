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
  IsWindowFocused

  constructor() {
    this.DocumentID = document.querySelector("html").dataset.id ||
      "random"+Math.random().toString(36) + Date.now().toString(36)
    this.WindowID = Math.random().toString(36) + '|' + Date.now()
    this.Curtain = document.querySelector('presentation-curtain')
    this.MainTimer = document.querySelector('presentation-timer#mainTimer')
    this.Channel = new BroadcastChannel(this.DocumentID)

    this.debouncedBroadcastWindowState = debounce(this.broadcastWindowState.bind(this), 120)
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
    this.onStateChange = this.onStateChange.bind(this)
    this.onFocus = this.onFocus.bind(this)
    this.onBlur = this.onBlur.bind(this)
    this.onMessage = this.onMessage.bind(this)

    this.setWindowFocused(document.hasFocus())
    window.addEventListener(navigationCompleteEventType, this.onNavigationComplete)
    this.Curtain?.addEventListener('change', this.onStateChange)
    this.MainTimer?.addEventListener('change', this.onStateChange)
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

    timer.setDuration(elapsed, duration)
    const running = state.IsTimerRunning
    if (running) {
      timer.start()
    } else {
      timer.pause()
    }
  }

  setMainCurtainState(state) {
    if (
      this.Curtain && typeof state?.IsCurtainDown === 'boolean' &&
      this.Curtain.IsDown() !== state.IsCurtainDown
    ) {
      if (state.IsCurtainDown) {
        this.Curtain.Down()
      } else {
        this.Curtain.Up()
      }
    }
  }

  broadcastWindowState(reloading = false) {
    if (this.IsWindowFocused !== true) {
      return
    }

    const state = new WindowState()
    state.DocumentID = this.DocumentID
    state.WindowID = this.WindowID
    state.SlideIndex = currentSlide
    state.ConcealedListItemCount = getCurrentConcealedListItems().length
    state.IsCurtainDown = Boolean(this.Curtain?.IsDown())
    state.IsReloading = reloading === true

    const timer = this.MainTimer
    const elapsed = Number(timer?.ElapsedDuration)
    const duration = Number(timer?.UtmostDuration)
    state.TimerElapsedDuration = elapsed
    state.TimerUtmostDuration = duration
    state.IsTimerRunning = timer.TimerRunning
    this.Channel.postMessage(state)
  }

  setWindowFocused(focused) {
    this.IsWindowFocused = focused
    document.body.classList.toggle('is-focused', focused)
  }

  onNavigationComplete() {
    this.debouncedBroadcastWindowState()
  }

  onStateChange(event) {
    if (event.target !== this.Curtain && event.target !== this.MainTimer) {
      return
    }
    this.broadcastWindowState()
  }

  onFocus() {
    this.setWindowFocused(true)
    // this.debouncedBroadcastWindowState()
  }

  onBlur() {
    this.setWindowFocused(false)
    // this.debouncedBroadcastWindowState()
  }

  onMessage(event) {
    const broadcast = event.data
    // if (broadcast.WindowID === this.WindowID) {
    //   return
    // }
    this.setWindowFocused(false)
    this.setMainCurtainState(broadcast)
    this.setMainTimerState(broadcast)
    if (!broadcast || broadcast.DocumentID !== this.DocumentID) {
      return
    }
    if (broadcast.IsReloading === true) {
      window.location.reload()
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
