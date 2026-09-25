// use as custom element <clock> anywhere on the page
(() => {
  const svgNamespace = 'http://www.w3.org/2000/svg'
  const element = (name, attributes = {}) => {
    const node = document.createElementNS(svgNamespace, name)
    for (const [attribute, value] of Object.entries(attributes)) {
      node.setAttribute(attribute, value)
    }
    return node
  }

  class MdcoachClock extends HTMLElement {
    constructor() {
      super()
      this.attachShadow({ mode: 'open' })
    }

    connectedCallback() {
      window.clearInterval(this.timer)

      const style = document.createElement('style')
      style.textContent = `
        :host {
          display: block;
          width: 100%;
          height: 100%;
        }

        svg {
          width: 100%;
          height: 100%;
        }

        .clock-face { stroke: #333; fill: white; }
        .minor { stroke: #999; stroke-width: 0.5; }
        .major { stroke: #333; stroke-width: 1; }
        .hour { stroke: #333; }
        .minute { stroke: #666; }
        .second, .second-counterweight { stroke: rgb(180, 0, 0); }
        .second-counterweight { stroke-width: 3; }
      `

      this.svg = element('svg', { viewBox: '-50 -50 100 100' })
      this.svg.append(element('circle', { class: 'clock-face', r: '48' }))

      for (let minute = 0; minute < 60; minute += 5) {
        this.svg.append(element('line', {
          class: 'major',
          y1: '35',
          y2: '45',
          transform: `rotate(${6 * minute})`,
        }))

        for (let offset = 1; offset <= 4; offset++) {
          this.svg.append(element('line', {
            class: 'minor',
            y1: '42',
            y2: '45',
            transform: `rotate(${6 * (minute + offset)})`,
          }))
        }
      }

      this.hourHand = element('line', { class: 'hour', y1: '2', y2: '-20' })
      this.minuteHand = element('line', { class: 'minute', y1: '4', y2: '-30' })
      this.secondHand = element('g')
      this.secondHand.append(
        element('line', { class: 'second', y1: '10', y2: '-38' }),
        element('line', { class: 'second-counterweight', y1: '10', y2: '2' }),
      )
      this.svg.append(this.hourHand, this.minuteHand, this.secondHand)
      this.shadowRoot.replaceChildren(style, this.svg)

      this.updateTime()
      this.timer = window.setInterval(() => this.updateTime(), 1000)
    }

    disconnectedCallback() {
      window.clearInterval(this.timer)
      this.timer = undefined
    }

    updateTime() {
      const time = new Date()
      const hours = time.getHours()
      const minutes = time.getMinutes()
      const seconds = time.getSeconds()

      this.hourHand.setAttribute('transform', `rotate(${30 * hours + minutes / 2})`)
      this.minuteHand.setAttribute('transform', `rotate(${6 * minutes + seconds / 10})`)
      this.secondHand.setAttribute('transform', `rotate(${6 * seconds})`)
    }
  }

  if (!customElements.get('clock')) {
    customElements.define('clock', MdcoachClock)
  }
})()
