export const scrollToPosition = (x, y, delay=15, limit=50) => {
  const deltaX = (window.scrollX - x) * 0.6
  const deltaY = (window.scrollY - y) * 0.6
  // console.log("deltaY", deltaY)
  if ((deltaY < 5 && deltaY > -5) || limit < 1) {
    window.scrollTo(x, y)
    return
  }
  window.scrollTo(x + deltaX, y + deltaY)
  setTimeout(() => scrollToPosition(x, y, delay, limit-1), delay)
}

const isVScrollNecessary = (element, tolerancePercent=20) => {
  const tolerance = element.clientHeight * tolerancePercent / 100
  // console.log("tole", tolerance, element)
  if (element.offsetTop > window.scrollY) {
    return element.offsetTop < window.scrollY + window.screen.height - tolerance
  } else {
    return element.offsetTop + element.clientHeight - tolerance > window.scrollY
  }
  // return true // not visible

  // const y = element.offsetTop - window.scrollY - element.height
  // // * tolerancePercent / 100
  // if (y < 0) return true                    // above window
  // if (y > window.screen.height) return true // below window
  // return false // visible
}

export const verticalScrollTo = (elementID) => {
  const element = document.getElementById(elementID)
  if (!element) throw new Error("element ID not found: "+elementID)
  // if (!isVScrollNecessary(element, 20)) return
  // const targetX = element.offsetLeft
  // const targetY = element.offsetTop
  // const deltaX = element.offsetLeft
  scrollToPosition(element.offsetLeft, element.offsetTop, 60)
}
