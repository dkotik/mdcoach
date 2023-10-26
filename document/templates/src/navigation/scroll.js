export const scrollToPosition = (x, y, delay=15, limit=15) => {
  const deltaX = Math.floor((window.scrollX - x) / 4)
  const deltaY = Math.floor((window.scrollY - y) / 4)
  if (!deltaY) return
  window.scrollTo(x + deltaX, y + deltaY)
  if (limit < 1) return
  setTimeout(() => scrollToPosition(x, y, delay, limit-1), delay)
}

export const scrollToID = (elementID) => {
  const element = document.getElementById(elementID)
  if (!element) throw new Error("element ID not found: "+elementID)
  // const targetX = element.offsetLeft
  // const targetY = element.offsetTop
  // const deltaX = element.offsetLeft
  scrollToPosition(element.offsetLeft, element.offsetTop)
}

export const isVerticalScrollNecessary = (elementID, tolerancePercent=20) => {
  const element = document.getElementById(elementID)
  if (!element) throw new Error("element ID not found: "+elementID)
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
