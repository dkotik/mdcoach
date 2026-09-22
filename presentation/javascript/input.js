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
    }
  }
)

document.addEventListener(
  'wheel',
  (event) => {
    if (event.deltaY >= 0) {
      next(event)
    } else {
      previous(event)
    }
  }
)
