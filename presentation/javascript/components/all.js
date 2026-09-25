const svgNamespace = 'http://www.w3.org/2000/svg'

const element = (name, attributes = {}) => {
  const node = document.createElementNS(svgNamespace, name)
  for (const [attribute, value] of Object.entries(attributes)) {
    node.setAttribute(attribute, value)
  }
  return node
}
