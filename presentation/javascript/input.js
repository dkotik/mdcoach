// ============== input.js =====================

let isCurtainVisible = false
const curtain = document.querySelector("body > aside.curtain")

const showCurtain = () => {
  curtain.style.display = "grid";
  isCurtainVisible = true
}

const hideCurtain = () => {
  curtain.style.display = "none";
  isCurtainVisible = false
}

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
    if (isCurtainVisible) {
      switch (event.code) {
        case 'Escape':
        case 'Period':
          hideCurtain()
      }
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
      case 'Period':
        showCurtain()
      case 'KeyR':
        window.location.reload()
    }
  }
)

document.addEventListener(
  'wheel',
  debounce((event) => {
    if (isCurtainVisible) {
      return
    }
    if (event.deltaY >= 0) {
      next(event)
    } else {
      previous(event)
    }
  }, 4)
)
