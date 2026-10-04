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

const documentID = document.querySelector("html").dataset.id ||
  "random"+Math.random().toString(36) + Date.now().toString(36)
const channel = new BroadcastChannel(documentID)
const windowID = documentID + '|' + Math.random().toString(36) + '|' + Date.now()
const windowState = new WindowState()
windowState.DocumentID = documentID
windowState.WindowID = windowID

const mainTimerSelector = 'presentation-timer#mainTimer'
const updateWindowState = (reload = false) => {
  windowState.SlideIndex = currentSlide
  windowState.ConcealedListItemCount = getCurrentConcealedListItems().length
  windowState.IsCurtainDown = Boolean(
    document.querySelector('presentation-curtain')?.hasAttribute('open'),
  )
  windowState.IsReloading = reload === true

  const timer = document.querySelector(mainTimerSelector)
  const duration = Number(timer?.duration)
  const elapsed = Number(timer?.elapsed) + (
    timer?.running
      ? performance.now() - Number(timer.startedAt)
      : 0
  )
  if (timer && Number.isFinite(duration) && duration > 0 && Number.isFinite(elapsed)) {
    const expired = timer.hasAttribute('expired') || elapsed >= duration
    windowState.TimerElapsedDuration = expired
      ? duration
      : Math.max(0, Math.min(duration, elapsed))
    windowState.TimerUtmostDuration = duration
    windowState.IsTimerRunning = timer.running && !expired
  } else {
    windowState.TimerElapsedDuration = 0
    windowState.TimerUtmostDuration = 0
    windowState.IsTimerRunning = false
  }
}
const setMainTimerState = (state) => {
  const timer = document.querySelector(mainTimerSelector)
  const duration = Number(state?.TimerUtmostDuration)
  const elapsed = Number(state?.TimerElapsedDuration)
  if (
    !timer || state?.DocumentID !== documentID ||
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
const setMainCurtainState = (state) => {
  const curtain = document.querySelector('presentation-curtain')
  if (curtain && typeof state?.IsCurtainDown === 'boolean') {
    curtain.toggleAttribute('open', state.IsCurtainDown)
  }
}
const broadcastWindowState = debounce((reload = false) => {
  updateWindowState(reload)
  channel.postMessage(windowState)
}, 120)
const setWindowFocused = (focused) => {
  document.body.classList.toggle('is-focused', focused)
}

setWindowFocused(document.hasFocus())

window.addEventListener(navigationCompleteEventType, () => {
  broadcastWindowState()
})
window.addEventListener(curtainToggleEventType, (event) => {
  if (event.target !== document.querySelector('presentation-curtain')) {
    return
  }
  broadcastWindowState()
})
window.addEventListener('focus', () => {
  setWindowFocused(true)
  broadcastWindowState()
})
window.addEventListener('blur', () => {
  setWindowFocused(false)
  broadcastWindowState()
})

channel.addEventListener('message', (event) => {
  const broadcast = event.data
  if (!broadcast || broadcast.DocumentID !== documentID) {
    return
  }
  if (broadcast.Reload === true) {
    window.location.reload()
    return
  }
  if (broadcast.WindowID === windowID) {
    return
  }
  setWindowFocused(false)

  setMainCurtainState(broadcast)
  setMainTimerState(broadcast)

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
})

broadcastWindowState()
