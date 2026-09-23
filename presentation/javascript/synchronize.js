// ============== synchronize.js =====================
const documentID = 'sync'
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

channel.addEventListener('message', (event) => {
  const broadcast = event.data
  if (broadcast.window === windowID) {
    return
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
