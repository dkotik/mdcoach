import { tick } from 'svelte'

function verticalScaleRatio(node) {
  const parentHeight = node.parentNode.parentNode.offsetHeight
  const margin = parentHeight * 0.1
  return () => parentHeight / (node.offsetHeight + margin)
}

export default function(node, condition) {
  let visible = condition
  let scaled = false

  const scale = async () => {
    if (scaled) {
      scaled = false
      node.style.transform = 'scale(1)'
      node.style.fontSize = '100%'
      await tick()
    } else {
      scaled = true
    }
    const ratioMeasure = verticalScaleRatio(node)
    let ratio = ratioMeasure()
    if (ratio > 1) return // no need

    // scale down fontSize first
    let fontSize = 100
    while (fontSize > 50) {
      fontSize -= 10
      node.style.fontSize = fontSize + '%'
      await tick()
      if (!visible) return
      ratio = ratioMeasure()
      if (ratio > 1) return
    }
    node.style.transform = 'scale(' + ratio + ')'
    // console.log("finished scaling!", node.style.fontSize)
  }
  if (visible) scale()

  const unlockScaling = async () => {
    if (visible) scale()
  }
  unlockScaling()
  window.addEventListener("resize", unlockScaling)
  return {
    update(condition) {
      visible = condition
      if (visible) scale()
    },

    destroy() {
      window.removeEventListener("resize", unlockScaling)
    }
  }
}
