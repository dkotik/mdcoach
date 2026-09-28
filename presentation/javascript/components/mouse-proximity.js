// use as <mouse-proximity distance="100"> to reveal slotted content near the pointer

const mouseProximityDefaultDistance = 100

class MouseProximity extends HTMLElement {
  static get observedAttributes() {
    return ['distance']
  }

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onPointerMove = this.onPointerMove.bind(this)
    this.onPointerOut = this.onPointerOut.bind(this)
    this.pointerX = undefined
    this.pointerY = undefined
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        display: inline-block;
      }

      .content {
        visibility: hidden;
      }

      :host([revealed]) .content {
        visibility: visible;
      }
    `

    this.content = document.createElement('span')
    this.content.className = 'content'
    this.content.append(document.createElement('slot'))
    this.shadowRoot.replaceChildren(style, this.content)
    this.updateState()
    window.addEventListener('pointermove', this.onPointerMove)
    window.addEventListener('pointerout', this.onPointerOut)
  }

  disconnectedCallback() {
    window.removeEventListener('pointermove', this.onPointerMove)
    window.removeEventListener('pointerout', this.onPointerOut)
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name !== 'distance' || oldValue === newValue || !this.isConnected) {
      return
    }

    this.updateRevealState()
  }

  getDistance() {
    const distance = Number(this.getAttribute('distance'))
    return Number.isFinite(distance) && distance >= 0
      ? distance
      : mouseProximityDefaultDistance
  }

  onPointerMove(event) {
    this.pointerX = event.clientX
    this.pointerY = event.clientY
    this.updateRevealState()
  }

  onPointerOut(event) {
    if (event.relatedTarget) {
      return
    }

    this.pointerX = undefined
    this.pointerY = undefined
    this.removeAttribute('revealed')
    this.updateState()
  }

  updateRevealState() {
    if (this.pointerX === undefined || this.pointerY === undefined) {
      return
    }

    const bounds = this.getBoundingClientRect()
    const horizontalDistance = Math.max(bounds.left - this.pointerX, 0, this.pointerX - bounds.right)
    const verticalDistance = Math.max(bounds.top - this.pointerY, 0, this.pointerY - bounds.bottom)
    this.toggleAttribute(
      'revealed',
      Math.hypot(horizontalDistance, verticalDistance) <= this.getDistance(),
    )
    this.updateState()
  }

  updateState() {
    this.setAttribute('aria-hidden', String(!this.hasAttribute('revealed')))
  }
}

if (!customElements.get('mouse-proximity')) {
  customElements.define('mouse-proximity', MouseProximity)
}
