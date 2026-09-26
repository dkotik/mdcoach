// ============== synchronize.js =====================
const documentID = document.querySelector("html").dataset.id ||
  "random"+Math.random().toString(36) + Date.now().toString(36)
const channel = new BroadcastChannel(documentID)
const windowID = documentID + '|' + Math.random().toString(36) + '|' + Date.now()

window.addEventListener(navigationCompleteEventType, (event) => {
  const navigationComplete = event.detail
  channel.postMessage({
    slideIndex: navigationComplete.slideIndex,
    concealedListItemCount: navigationComplete.concealedListItemCount,
    window: windowID
  });
})

let isWindowFocused = true
channel.addEventListener('message', (event) => {
  const broadcast = event.data
  if (broadcast.window === windowID) {
    if (!isWindowFocused) {
      isWindowFocused = true
      document.documentElement.querySelector("body").classList.add("is-focused")
    }
    return
  }

  if (isWindowFocused) {
    isWindowFocused = false
    document.documentElement.querySelector("body").classList.remove("is-focused")
  }

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
