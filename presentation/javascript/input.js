// ============== input.js =====================
const nextEventType = "nextListItemOrSlide"
const next = (event) => {
  event.preventDefault()
  window.dispatchEvent(new CustomEvent(nextEventType, {
    bubbles: true,
    cancelable: true,
    detail: event
  }))
}

const previousEventType = "previousSlide"
const previous = (event) => {
  event.preventDefault()
  window.dispatchEvent(new CustomEvent(previousEventType, {
    bubbles: true,
    cancelable: true,
    detail: event
  }))
}

document.addEventListener(
  'keydown',
  (event) => {
    if (document.querySelector('presentation-curtain[open]')) {
      return
    }

    switch (event.code) {
      case 'ArrowRight':
      case 'ArrowDown':
      case 'PageDown':
      case 'Space':
      case 'KeyJ':
        next(event)
        return
      case 'ArrowLeft':
      case 'ArrowUp':
      case 'PageUp':
      case 'Backspace':
      case 'KeyK':
        previous(event)
        return
      case 'KeyC':
        window.open(window.location.href, '_blank')
        return
      case 'KeyR':
        window.location.reload()
    }
  }
)

document.getElementById("slideJump").addEventListener(
  keyStrokeCompleteEventType,
  (event) => {
    const requestedSlide = Number.parseInt(event.detail.replace(/\s+/g, ''), 10)
    if (!Number.isInteger(requestedSlide)) {
      return
    }

    const slideNumber = Math.min(
      Math.max(requestedSlide, 1),
      finalSlideIndex + 1,
    )
    navigate(slideNumber - 1)
    concealedListItems = getCurrentConcealedListItems()
    dispatchNavigationCompleteEvent(concealedListItems.length)
  }
)
