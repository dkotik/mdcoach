// use as <presentation-timer duration="60">; duration is in seconds

class PresentationTimer extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.elapsed = 0
    this.running = false
    this.frame = undefined
    this.startedAt = 0
    this.onClick = this.onClick.bind(this)
    this.onKeyDown = this.onKeyDown.bind(this)
    this.tick = this.tick.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        cursor: pointer;
        display: block;
        height: 100%;
        width: 100%;
      }

      :host(:focus-visible) {
        outline: 2px solid currentColor;
        outline-offset: 2px;
      }

      svg {
        height: 100%;
        width: 100%;
      }

      .track {
        fill: none;
        stroke: var(--color-body-subtext, #aaa);
        stroke-opacity: 0.35;
        stroke-width: 16;
      }

      .progress {
        fill: none;
        stroke: var(--color-marker-background, #006eff);
        stroke-linecap: round;
        stroke-width: 16;
        transition: stroke-opacity 150ms ease;
      }

      :host([paused]) .progress {
        stroke-opacity: 0.3;
      }

      @media (prefers-reduced-motion: reduce) {
        .progress {
          stroke-linecap: butt;
        }
      }
    `

    this.duration = this.getDuration()
    this.circumference = 2 * Math.PI * 42
    this.svg = element('svg', { viewBox: '0 0 100 100', 'aria-hidden': 'true' })
    this.trackCircle = element('circle', {
      class: 'track',
      cx: '50',
      cy: '50',
      r: '36',
    })
    this.progressCircle = element('circle', {
      class: 'progress',
      cx: '50',
      cy: '50',
      r: '36',
      transform: 'rotate(-90 50 50)',
      'stroke-dasharray': String(this.circumference),
      'stroke-dashoffset': String(this.circumference),
    })
    this.svg.append(this.trackCircle, this.progressCircle)
    this.shadowRoot.replaceChildren(style, this.svg)

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'button')
    }
    if (!this.hasAttribute('tabindex')) {
      this.setAttribute('tabindex', '0')
    }

    this.setAttribute('aria-label', `${this.duration / 1000} second timer`)
    this.addEventListener('click', this.onClick)
    this.addEventListener('keydown', this.onKeyDown)
    this.updateProgress(this.elapsed)
    this.updateState()
  }

  disconnectedCallback() {
    this.pause()
    this.removeEventListener('click', this.onClick)
    this.removeEventListener('keydown', this.onKeyDown)
  }

  getDuration() {
    const seconds = Number(this.getAttribute('duration'))
    return Number.isFinite(seconds) && seconds > 0 ? seconds * 1000 : 60_000
  }

  onClick() {
    if (this.running) {
      this.pause()
    } else {
      this.start()
    }
  }

  onKeyDown(event) {
    if ((event.code !== 'Enter' && event.code !== 'Space') || event.repeat) {
      return
    }

    event.preventDefault()
    this.onClick()
  }

  start() {
    if (this.elapsed >= this.duration) {
      this.elapsed = 0
    }

    this.startedAt = performance.now()
    this.running = true
    this.updateProgress(this.elapsed)
    this.updateState()
    this.frame = window.requestAnimationFrame(this.tick)
  }

  pause() {
    if (!this.running) {
      return
    }

    this.elapsed = Math.min(
      this.duration,
      this.elapsed + performance.now() - this.startedAt,
    )
    this.running = false
    window.cancelAnimationFrame(this.frame)
    this.frame = undefined
    this.updateProgress(this.elapsed)
    this.updateState()
  }

  tick(timestamp) {
    if (!this.running) {
      return
    }

    const elapsed = this.elapsed + timestamp - this.startedAt
    if (elapsed >= this.duration) {
      this.elapsed = this.duration
      this.running = false
      this.frame = undefined
      this.updateProgress(this.elapsed)
      this.updateState()
      return
    }

    this.updateProgress(elapsed)
    this.frame = window.requestAnimationFrame(this.tick)
  }

  updateProgress(elapsed) {
    const progress = Math.min(elapsed / this.duration, 1)
    this.progressCircle.setAttribute(
      'stroke-dashoffset',
      String(this.circumference * (1 - progress)),
    )
  }

  updateState() {
    this.toggleAttribute('paused', !this.running)
    this.setAttribute('aria-pressed', String(this.running))
    this.setAttribute(
      'aria-label',
      `${this.running ? 'Pause' : 'Start'} ${this.duration / 1000} second timer`,
    )
  }
}

if (!customElements.get('presentation-timer')) {
  customElements.define('presentation-timer', PresentationTimer)
}
