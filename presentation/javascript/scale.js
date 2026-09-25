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

/*
Some legacy functionality:

const scale = async () => {
  if (scaled) {
    scaled = false // reset previously set scale
    node.style.marginTop = '0' // reset top margin
    node.style.transform = 'scale(1)'
    node.style.fontSize = '100%'
    // await tick()
  } else {
    scaled = true
  }
  const ratioMeasure = verticalScaleRatio(node)
  let ratio = await ratioMeasure()
  if (ratio > 1) return // no need

  // scale down fontSize first
  let fontSize = 100
  while (fontSize > 50) {
    fontSize -= 10
    node.style.fontSize = fontSize + '%'
    // await tick() // TODO: tick does not seem to work, replace with timer?
    if (!visible) return
    ratio = await ratioMeasure()
    if (ratio > 1) return
  }
  node.style.transform = 'scale(' + ratio + ')'
  // await delayPromise(100)
  // console.log("finished scaling!", node.style.fontSize)
}

const scaleThenCenter = async () => {
  if (!visible) return
  await scale()
  const slideElement = node?.parentNode
  if (!slideElement) return
  await tick()
  // await delayPromise(120)
  const parentHeight = slideElement.parentNode?.offsetHeight || window.height
  const gap = parentHeight - node.offsetHeight
  node.style.marginTop = Math.floor(gap*0.5) + "px"
  // console.log(node.parentNode, (gap/2) + "px gap")
  // console.log((gap/2) + "px gap")
}

*/
