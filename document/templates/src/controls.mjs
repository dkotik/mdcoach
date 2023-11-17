import { tick } from 'svelte'

export function debounce(callback, delay=90) {
  let timer
  return (...args) => {
    clearTimeout(timer)
    timer = setTimeout(() => callback(...args), delay)
  }
}

const allRevealedItems = new Map()
export function revealedListItems(node, data) {
  // query='section.active > article > ul > li:not(.revealed)'
  let target = null
  let currentSlide = 1

  const gatherConcealedListItems = (event) => {
    target = node.querySelector(':scope > section.active > article')
    if (!target) return
    if (!target.listItems) target.listItems = Array.from(target.querySelectorAll(':scope > ul > li'))
    if (target.listItems.length === 0) return
    // reveal items that should be visible based on previous depth
    for (var i = 0; i < allRevealedItems.get(currentSlide); i++) {
      target.listItems[i]?.classList.add('revealed')
    }
  }

  const nextConcealedListItem = (event) => {
    let nextIndex = allRevealedItems.get(currentSlide) || 0
    const nextConcealed = target?.listItems[nextIndex]

    if (nextConcealed) {
      nextIndex++
      allRevealedItems.set(currentSlide, nextIndex)
      nextConcealed.classList.add('revealed')
      return node.dispatchEvent(new CustomEvent('nextListItem', {
        bubbles: true, cancelable: true,
        detail: { slide: currentSlide, listItem: nextIndex }
      }))
    }

    return node.dispatchEvent(new CustomEvent('nextListItem', {
      bubbles: true, cancelable: true,
      detail: { slide: currentSlide+1, listItem: 0 }
    }))
  }

  node.addEventListener('next', nextConcealedListItem)
  node.addEventListener('nextStop', gatherConcealedListItems)
  node.addEventListener('previousStop', gatherConcealedListItems)
  node.addEventListener('scrollStop', gatherConcealedListItems)
  node.addEventListener('jump', gatherConcealedListItems)
  tick().then(gatherConcealedListItems)
  const gatherWithDelay = debounce(gatherConcealedListItems, 200)
  return {
    update([slide, listItem]) {
      currentSlide = slide
      gatherWithDelay()
      const revealedIndex = allRevealedItems.get(slide) || 0
      if (revealedIndex < listItem) {
        allRevealedItems.set(slide, listItem)
        // gatherConcealedListItems()
        return
      }
      // gatherConcealedListItems()
      // if (target) {
      //   for (var i = 0; i < allRevealedItems.get(currentSlide); i++) {
      //     target.listItems[i].classList.add('revealed')
      //   }
      // }
      node.dispatchEvent(new CustomEvent('nextListItem', {
        bubbles: true, cancelable: true,
        detail: { slide: slide, listItem: revealedIndex }
      }))

      // if (!allRevealedItems.has(slide) || listItem > allRevealedItems.get(slide)) {
      //   allRevealedItems.set(slide, listItem)
      //   tick().then(gatherConcealedListItems)
      //   // for (var i = 0; i < listItem; i++) {
      //   //   target.listItems[i].classList.add('revealed')
      //   // }
      //
      // }
      // gatherConcealedListItems()
      // setTimeout(gatherConcealedListItems, 600)
      // // reveal items that should be visible based on previous depth
      // for (var i = 0; i < listItem; i++) {
      //   target.listItems[i].classList.add('revealed')
      // }
      // console.log("updated to", slide, listItem)
    },

    destroy() {
      node.removeEventListener('jump', gatherConcealedListItems)
      node.removeEventListener('scrollStop', gatherConcealedListItems)
      node.removeEventListener('previousStop', gatherConcealedListItems)
      node.removeEventListener('nextStop', gatherConcealedListItems)
      node.removeEventListener('next', nextConcealedListItem)
    }
  }
}

export function keyboardNavigation(node) {
  const handleKeyDownEvent = (event) => {
    switch (event.code) {
      case 'ArrowRight':
      case 'ArrowDown':
      case 'PageDown':
      case 'Space':
      case 'KeyJ':
        event.preventDefault()
        node.dispatchEvent(new CustomEvent('next', {
          bubbles: true,
          cancelable: true,
          detail: event
        }))
        return
      case 'ArrowLeft':
      case 'ArrowUp':
      case 'PageUp':
      case 'Backspace':
      case 'KeyK':
        event.preventDefault()
        node.dispatchEvent(new CustomEvent('previous', {
          bubbles: true,
          cancelable: true,
          detail: event
        }))
        return
      case 'KeyC':
        window.open(window.location.href, '_blank')
        return
    }
  }

  let rememberedDigits = ''
  const forgetDigits = debounce(() => rememberedDigits = '', 3000)
  const rememberDigit = (digit) => {
    rememberedDigits += digit
    forgetDigits()
  }
  const emitJumpEvent = (node) => {
    if (!rememberedDigits) {
      return false
    }
    event.preventDefault()
    node.dispatchEvent(new CustomEvent('jump', {
      bubbles: true,
      cancelable: true,
      detail: parseInt(rememberedDigits),
    }))
    rememberedDigits = ''
    return true
  }

  const handleKeyUpEvent = (event) => {
    switch (event.code) {
      case 'ArrowRight':
      case 'ArrowDown':
      case 'PageDown':
      case 'Space':
      case 'KeyJ':
        if (emitJumpEvent(node)) return
        node.dispatchEvent(new CustomEvent('nextStop', {
          bubbles: true,
          cancelable: true,
          detail: parseInt(rememberedDigits),
        }))
        return
      case 'ArrowLeft':
      case 'ArrowUp':
      case 'PageUp':
      case 'Backspace':
      case 'KeyK':
        node.dispatchEvent(new CustomEvent('previousStop', {
          bubbles: true,
          cancelable: true,
          detail: parseInt(rememberedDigits),
        }))
        return
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
      case 'KeyG': emitJumpEvent(node); return
    }
  }

  document.addEventListener('keydown', handleKeyDownEvent)
  document.addEventListener('keyup', handleKeyUpEvent)
  return {
    destroy() {
      document.removeEventListener('keydown', handleKeyDownEvent)
      document.removeEventListener('keyup', handleKeyUpEvent)
    }
  }
}

export function wheelNavigation(node) {
  const handleWheelEvent = (event) => {
    // event.preventDefault()
    if (event.deltaY >= 0) {
      node.dispatchEvent(new CustomEvent('next', {
        bubbles: true,
        cancelable: true,
        detail: event
      }))
    } else {
      node.dispatchEvent(new CustomEvent('previous', {
        bubbles: true,
        cancelable: true,
        detail: event
      }))
    }
  }

  document.addEventListener('wheel', handleWheelEvent)
  return {
    destroy() {
      document.removeEventListener('wheel', handleWheelEvent)
    }
  }
}

export function scrollStop(node, delay=150) {
  const detectStop = debounce((event) => {
    if (!document.hasFocus()) {
      // console.log("dropping scroll stop, cuz not focused")
      return
    }
    node.dispatchEvent(new CustomEvent('scrollStop', {}))
  }, delay)

  node.addEventListener('scroll', detectStop)
  return {
    destroy() {
      node.removeEventListener('scroll', detectStop)
    }
  }
}
