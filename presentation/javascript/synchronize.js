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

const broadcastWindowState = () => channel.postMessage({
  slideIndex: navigationState.slideIndex,
  concealedListItemCount: navigationState.concealedListItemCount,
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
