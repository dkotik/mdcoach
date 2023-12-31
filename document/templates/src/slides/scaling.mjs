import { tick } from 'svelte'

async function delayPromise(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function verticalScaleRatio(node) {
  const parentHeight = node?.parentNode?.parentNode?.offsetHeight || window.height
  const margin = parentHeight * 0.1
  return async () => {
    await delayPromise(20)
    return parentHeight / (node.offsetHeight + margin)
  }
}

export default function(node, condition) {
  let visible = condition
  let scaled = false

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

  const unlockScaling = async () => {
    scaled = false
    await scaleThenCenter()
  }
  unlockScaling()
  window.addEventListener("resize", unlockScaling)
  node.style.marginTop = '0' // reset top margin
  return {
    update(condition) {
      visible = condition
      scaleThenCenter()
    },

    destroy() {
      window.removeEventListener("resize", unlockScaling)
    }
  }
}
