// ============== scale.js =====================

function scaleElementDown(element) {
  const attribute = element.getAttribute("data-scale")
  const current = Number.parseInt(attribute ?? "100", 10)
  if (current < 20) {
    return false
  }
  element.setAttribute("data-scale", String(current - 10))
  return true
}

let topEdge = 20
function isTooLarge(element) {
  return (element.getBoundingClientRect().height + topEdge) > window.innerHeight
}

function scaleDown() {
  let scaling = false
  for (const slide of slides) {
    const element = slide.querySelector(":scope > .grid > .content")
    if (isTooLarge(element)) {
      scaling = scaleElementDown(element)
    }
  }

  if (scaling) {
    window.setTimeout(scaleDown, 150)
  }
}

window.setTimeout(() => {
  // topEdge = slides[0]?.getBoundingClientRect().top ?? 20
  // console.log("topEdge", topEdge)
  topEdge = parseFloat(window.getComputedStyle(slides[0].querySelector(":scope > .grid")).paddingTop)
  scaleDown()
}, 150)

window.addEventListener("resize", debounce(() => {
  for (const slide of slides) {
    const element = slide.querySelector(":scope > .grid > .content")
    element.setAttribute("data-scale", "100")
  }
  window.setTimeout(scaleDown, 150)
}, 200))
