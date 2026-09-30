// ============== synchronize.js =====================
const documentID = document.querySelector("html").dataset.id ||
  "random"+Math.random().toString(36) + Date.now().toString(36)
const channel = new BroadcastChannel(documentID)
const windowID = documentID + '|' + Math.random().toString(36) + '|' + Date.now()

let navigationState = {
  slideIndex: currentSlide,
  concealedListItemCount: getCurrentConcealedListItems().length,
}
let isWindowFocused = document.hasFocus()
const bodyElement = document.body
bodyElement.classList.toggle('is-focused', isWindowFocused)

const mainTimerSelector = 'presentation-timer#mainTimer'
const getMainTimerState = () => {
  const timer = document.querySelector(mainTimerSelector)
  if (!timer) {
    return null
  }

  const now = performance.now()
  const duration = Number(timer.duration)
  const elapsed = Number(timer.elapsed) + (
    timer.running
      ? now - Number(timer.startedAt)
      : 0
  )
  const expired = timer.hasAttribute('expired') || elapsed >= duration

  return {
    timerID: timer.id,
    duration,
    remainingDuration: expired ? 0 : Math.max(0, duration - elapsed),
    running: timer.running && !expired,
    expired,
  }
}
const setMainTimerState = (state) => {
  const timer = document.querySelector(mainTimerSelector)
  const duration = Number(state?.duration)
  const remainingDuration = Number(state?.remainingDuration)
  if (
    !timer || state?.timerID !== timer.id ||
    !Number.isFinite(duration) || duration <= 0 ||
    !Number.isFinite(remainingDuration) || remainingDuration < 0 ||
    remainingDuration > duration || typeof state.running !== 'boolean' ||
    typeof state.expired !== 'boolean'
  ) {
    return
  }

  const expired = state.expired || remainingDuration === 0
  const boundedRemainingDuration = expired ? 0 : remainingDuration
  const running = state.running && boundedRemainingDuration > 0 && !expired
  timer.pause()
  timer.duration = duration
  timer.elapsed = duration - boundedRemainingDuration
  timer.setExpired(expired)
  timer.updateProgress(timer.elapsed)
  timer.updateState()
  if (running) {
    timer.start()
  }
}
const broadcastWindowState = () => channel.postMessage({
  slideIndex: navigationState.slideIndex,
  concealedListItemCount: navigationState.concealedListItemCount,
  timerState: getMainTimerState(),
  window: windowID,
})
const setWindowFocused = (focused) => {
  isWindowFocused = focused
  bodyElement.classList.toggle('is-focused', focused)
}

window.addEventListener(navigationCompleteEventType, (event) => {
  navigationState = event.detail
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
  if (broadcast.window === windowID) {
    setWindowFocused(broadcast.focused === true)
    return
  }
  setWindowFocused(false)

  setMainTimerState(broadcast.timerState)

  if (broadcast.slideIndex !== currentSlide) {
    navigate(broadcast.slideIndex)
  }
  const concealedListItems = getCurrentConcealedListItems()
  let revealCount = concealedListItems.length - broadcast.concealedListItemCount
  for (const element of concealedListItems) {
    if (revealCount === 0) {
      return
    }
    element.classList.add("is-revealed")
    revealCount--
  }
});

broadcastWindowState()
