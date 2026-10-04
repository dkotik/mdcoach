// use as <presentation-timer duration="60">; duration is in seconds

const presentationTimerTickEventType = 'presentation-timer-tick'

class PresentationTimer extends HTMLElement {
  ElapsedDuration
  UtmostDuration
  TimerRunning

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.ElapsedDuration = 0
    this.UtmostDuration = 0
    this.TimerRunning = false
    this.frame = undefined
    this.lastFrameTimestamp = undefined
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

      .expired-indicator {
        fill: #ffd641;
        opacity: 0;
      }

      .expired-indicator-square {
        fill: #ff4500;
        opacity: 0;
      }

      :host([expired]) .expired-indicator,
      :host([expired]) .expired-indicator-square {
        animation: timer-expired-pulse 800ms ease-in-out infinite;
        opacity: 1;
      }

      @keyframes timer-expired-pulse {
        0%, 100% {
          opacity: 0.15;
        }
        50% {
          opacity: 1;
        }
      }

      @media (prefers-reduced-motion: reduce) {
        .progress {
          stroke-linecap: butt;
        }

        :host([expired]) .expired-indicator,
        :host([expired]) .expired-indicator-square {
          animation: none;
          opacity: 1;
        }
      }
    `

    this.UtmostDuration = this.getDuration()
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
    this.expiredIndicatorBackground = element('circle', {
      class: 'expired-indicator',
      cx: '50',
      cy: '50',
      r: '34',
    })
    this.expiredIndicator = element('rect', {
      class: 'expired-indicator-square',
      x: '35',
      y: '35',
      width: '30',
      height: '30',
    })
    this.svg.append(
      this.trackCircle,
      this.progressCircle,
      this.expiredIndicatorBackground,
      this.expiredIndicator,
    )

    this.shadowRoot.replaceChildren(style, this.svg)

    if (!this.hasAttribute('role')) {
      this.setAttribute('role', 'button')
    }
    if (!this.hasAttribute('tabindex')) {
      this.setAttribute('tabindex', '0')
    }

    this.setAttribute('aria-label', `${this.UtmostDuration / 1000} second timer`)
    this.addEventListener('click', this.onClick)
    this.addEventListener('keydown', this.onKeyDown)
    this.restoreExpirationState()
    this.updateProgress(this.ElapsedDuration)
    this.updateState()
  }

  disconnectedCallback() {
    this.pause()
    this.removeEventListener('click', this.onClick)
    this.removeEventListener('keydown', this.onKeyDown)
  }

  getExpirationStorageKey() {
    const presentationID = document.documentElement?.dataset.id
    if (!presentationID) {
      return null
    }

    return `${presentationID}:timer:${this.id || 'default'}:expired`
  }

  restoreExpirationState() {
    const key = this.getExpirationStorageKey()
    if (!key) {
      return
    }

    try {
      this.toggleAttribute('expired', window.localStorage.getItem(key) === 'true')
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

  setExpired(isExpired) {
    this.toggleAttribute('expired', isExpired)

    const key = this.getExpirationStorageKey()
    if (!key) {
      return
    }

    try {
      if (isExpired) {
        window.localStorage.setItem(key, 'true')
      } else {
        window.localStorage.removeItem(key)
      }
    } catch {
      // Local storage may be unavailable in restricted browsing contexts.
    }
  }

  getDuration() {
    const seconds = Number(this.getAttribute('duration'))
    return Number.isFinite(seconds) && seconds > 0 ? seconds * 1000 : 60_000
  }

  onClick() {
    if (this.TimerRunning) {
      this.pause()
    } else if (this.hasAttribute('expired')) {
      this.setExpired(false)
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
    if (this.ElapsedDuration >= this.UtmostDuration) {
      this.ElapsedDuration = 0
    }
    this.setExpired(false)

    this.lastFrameTimestamp = performance.now()
    this.TimerRunning = true
    this.updateProgress(this.ElapsedDuration)
    this.updateState()
    this.frame = window.requestAnimationFrame(this.tick)
  }

  pause() {
    if (!this.TimerRunning) {
      return
    }

    const previousElapsed = this.ElapsedDuration
    this.ElapsedDuration = Math.min(
      this.UtmostDuration,
      this.ElapsedDuration + Math.max(0, performance.now() - this.lastFrameTimestamp),
    )
    this.dispatchElapsedSeconds(previousElapsed, this.ElapsedDuration)
    this.TimerRunning = false
    window.cancelAnimationFrame(this.frame)
    this.frame = undefined
    this.lastFrameTimestamp = undefined
    this.updateProgress(this.ElapsedDuration)
    this.updateState()
  }

  tick(timestamp) {
    if (!this.TimerRunning) {
      return
    }

    const previousElapsed = this.ElapsedDuration
    this.ElapsedDuration = Math.min(
      this.UtmostDuration,
      this.ElapsedDuration + Math.max(0, timestamp - this.lastFrameTimestamp),
    )
    this.lastFrameTimestamp = timestamp
    this.dispatchElapsedSeconds(previousElapsed, this.ElapsedDuration)
    if (this.ElapsedDuration >= this.UtmostDuration) {
      if (this.UtmostDuration % 1000 !== 0) {
        this.dispatchTimerTick(this.UtmostDuration)
      }
      this.TimerRunning = false
      this.frame = undefined
      this.lastFrameTimestamp = undefined
      this.setExpired(true)
      this.updateProgress(this.ElapsedDuration)
      this.updateState()
      return
    }

    this.updateProgress(this.ElapsedDuration)
    this.frame = window.requestAnimationFrame(this.tick)
  }

  dispatchElapsedSeconds(previousElapsed, elapsed) {
    const previousSecond = Math.floor(previousElapsed / 1000)
    const elapsedSecond = Math.floor(elapsed / 1000)
    for (let second = previousSecond + 1; second <= elapsedSecond; second++) {
      this.dispatchTimerTick(Math.min(second * 1000, this.UtmostDuration))
    }
  }

  dispatchTimerTick(elapsed) {
    this.dispatchEvent(new CustomEvent(presentationTimerTickEventType, {
      bubbles: true,
      composed: true,
      detail: {
        totalDuration: this.UtmostDuration,
        remainingDuration: Math.max(0, this.UtmostDuration - elapsed),
      },
    }))
  }

  updateProgress(elapsed) {
    const progress = Math.min(elapsed / this.UtmostDuration, 1)
    this.progressCircle.setAttribute(
      'stroke-dashoffset',
      String(this.circumference * (1 - progress)),
    )
  }

  updateState() {
    this.toggleAttribute('paused', !this.TimerRunning)
    this.setAttribute('aria-pressed', String(this.TimerRunning))
    this.setAttribute(
      'aria-label',
      `${this.TimerRunning ? 'Pause' : 'Start'} ${this.UtmostDuration / 1000} second timer`,
    )
  }
}

if (!customElements.get('presentation-timer')) {
  customElements.define('presentation-timer', PresentationTimer)
}
