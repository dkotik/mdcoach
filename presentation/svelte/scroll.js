let scrollTimeout = null
export const scrollToPosition = (x, y, delay=15, limit=50) => {
  const deltaX = (window.scrollX - x) * 0.6
  const deltaY = (window.scrollY - y) * 0.6
  // console.log("deltaY", deltaY)
  if ((deltaY < 5 && deltaY > -5) || limit < 1) {
    window.scrollTo(x, y)
    return
  }
  window.scrollTo(x + deltaX, y + deltaY)
  if (scrollTimeout) clearTimeout(scrollTimeout)
  scrollTimeout = setTimeout(() => scrollToPosition(x, y, delay, limit-1), delay)
}

export const isVerticalScrollNecessary = (elementID) => {
  const element = document.getElementById(elementID)
  if (!element) throw new Error("element ID not found: "+elementID)
  // console.log("tole", tolerance, element)
  if (window.scrollY + window.screen.height <= element.offsetTop + element.clientHeight * 1.2) {
    // return element.offsetTop + element.clientHeight - tolerance < window.scrollY + window.screen.height
    // console.log("above", window.scrollY, element.offsetTop)
    return true
  } else if (window.scrollY > element.offsetTop + element.clientHeight * 0.2) {
    // console.log("below", window.scrollY, element.offsetTop - element.clientHeight)
    return true
    // return element.offsetTop < window.scrollY + window.screen.height - tolerance
  }
  return false
}
