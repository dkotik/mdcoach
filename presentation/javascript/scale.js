// ============== scale.js =====================

function scaleElementDown(element) {
  const attribute = element.getAttribute("data-scale")
  const current = Number.parseInt(attribute ?? "100", 10)
  if (current < 20) {
    return false
  }
  element.setAttribute("data-scale", String(value - 10))
  return true
}

function isTooLarge(element) {
  return element.getBoundingClientRect().height * 1.1 > window.innerHeight
}

function scaleDown() {
  const scaling = false
  for (const element of slide) {
    if (isTooLarge(element)) {
      scaling = scaleElementDown(element)
    }
  }

  if (scaling) {
    window.setTimeout(scaleDown, 50)
  }
}

window.setTimeout(scaleDown, 150)
