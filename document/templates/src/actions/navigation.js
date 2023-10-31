export default function(node, wheel=false) {
  const handleKeyDownEvent = (event) => {
    switch (event.code) {
      case 'ArrowRight':
      case 'ArrowDown':
      case 'PageDown':
      case 'Space':
      case 'KeyJ':
        event.preventDefault()
        node.dispatchEvent(new CustomEvent('next', { detail: event }))
        return
      case 'ArrowLeft':
      case 'ArrowUp':
      case 'PageUp':
      case 'Backspace':
      case 'KeyK':
        event.preventDefault()
        node.dispatchEvent(new CustomEvent('previous', { detail: event }))
        return
      case 'KeyC':
        window.open(window.location.href, '_blank')
        return
    }
  }

  let rememberedDigits = ''
  let rememberedDigitsTimeout
  const rememberDigit = (digit) => {
    rememberedDigits += digit
    clearTimeout(rememberedDigitsTimeout)
    rememberedDigitsTimeout = setTimeout(() => {
      rememberedDigits = ''
    }, 3000)
  }
  const handleKeyUpEvent = (event) => {
    switch (event.code) {
      case 'Digit0': rememberDigit('0'); return
      case 'Digit1': rememberDigit('1'); return
      case 'Digit2': rememberDigit('2'); return
      case 'Digit3': rememberDigit('3'); return
      case 'Digit4': rememberDigit('4'); return
      case 'Digit5': rememberDigit('5'); return
      case 'Digit6': rememberDigit('6'); return
      case 'Digit7': rememberDigit('7'); return
      case 'Digit8': rememberDigit('8'); return
      case 'Digit9': rememberDigit('9'); return
      case 'KeyG':
        node.dispatchEvent(new CustomEvent('jump', {
          detail: parseInt(rememberedDigits),
        }))
      default:
        rememberedDigits = ''
        clearTimeout(rememberedDigitsTimeout)
    }
  }

  const handleWheelEvent = (event) => {
    event.preventDefault()
    if (event.deltaY >= 0) {
      node.dispatchEvent(new CustomEvent('next', { detail: event }))
    } else {
      node.dispatchEvent(new CustomEvent('previous', { detail: event }))
    }
  }

  document.addEventListener('keydown', handleKeyDownEvent)
  document.addEventListener('keyup', handleKeyUpEvent)
  if (wheel) document.addEventListener('wheel', handleWheelEvent)
  return {
    destroy() {
      document.removeEventListener('keydown', handleKeyDownEvent)
      document.removeEventListener('keyup', handleKeyUpEvent)
      if (wheel) document.removeEventListener('wheel', handleWheelEvent)
    }
  }
}
