export default function(node, delay = 150) {
  let isScrolling
  const detectStop = (event) => {
    clearTimeout(isScrolling)
    isScrolling = setTimeout(() => {
      node.dispatchEvent(new CustomEvent('scrollStop', {}))
    }, delay)
  }

  node.addEventListener('scroll', detectStop)
  return {
    destroy() {
      node.removeEventListener('scroll', detectStop)
    }
  }
}
