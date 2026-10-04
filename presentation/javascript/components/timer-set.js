// use as <timer-set duration="300">; duration is in seconds
// emits 'presentation-timer-change' with the selected duration in seconds

class TimerSet extends HTMLElement {
  static get observedAttributes() {
    return ['data-timer-id']
  }

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onClick = this.onClick.bind(this)
    this.onInputChange = this.onInputChange.bind(this)
    this.onInputKeyDown = this.onInputKeyDown.bind(this)
    this.onInputBlur = this.onInputBlur.bind(this)
    this.onTimerTick = this.onTimerTick.bind(this)
    this.defaultDuration = 60
    this.duration = this.defaultDuration
    this.timerElement = null
  }

  connectedCallback() {
    this.updateTimerTickSubscription()
    this.defaultDuration = this.getDurationAttribute()
    this.duration = this.defaultDuration

    const style = document.createElement('style')
    style.textContent = `
      :host {
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-body-subtext, #aaa);
        color: var(--color-body-text);
        display: inline-flex;
        max-width: 100%;
        margin-top: 0.6rem;
        font-size: 1rem;
      }

      .controls {
        align-items: center;
        display: flex;
        gap: 0.5rem;
      }


      button {
        color: var(--color-body-text);
        border: 1px solid transparent;
        background: transparent;
        border-radius: 0.3rem;
        cursor: pointer;
        font: inherit;
        line-height: 1.2;
        padding: 0.3rem 0.9rem;
        white-space: nowrap;
      }

      button:hover {
        color: var(--color-marker-background, #006eff);
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-marker-background, #006eff);
      }

      button:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      .duration-input {
        color: var(--color-body-text);
        border: 1px solid transparent;
        background: transparent;
        border-radius: 0.35rem;
        box-sizing: content-box;
        font: inherit;
        font-variant-numeric: tabular-nums;
        font-weight: 700;
        padding: 0.45rem 0.6rem;
        text-align: center;
        width: 5ch;

        &:hover {
          background: var(--color-menu-background, #eee);
          border: 1px solid var(--color-body-subtext, #aaa);
        }
      }

      .duration-input:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      @media (max-width: 30rem) {
        .controls {
          gap: 0.25rem;
        }

        button {
          font-size: 0.8rem;
          padding: 0.3rem;
        }
      }
    `

    this.input = document.createElement('input')
    this.input.className = 'duration-input'
    this.input.type = 'text'
    this.input.setAttribute('aria-label', 'Timer duration in minutes and seconds')
    this.input.setAttribute('autocomplete', 'off')
    this.input.addEventListener('change', this.onInputChange)
    this.input.addEventListener('keydown', this.onInputKeyDown)
    this.input.addEventListener('blur', this.onInputBlur)

    const controls = document.createElement('div')
    controls.className = 'controls'
    controls.append(
      this.createButton('Reset timer', '↻', 'reset'),
      this.createButton('Subtract one minute', '−1', -61),
      this.createButton('Add 1 minute', '+1', 59),
      this.createButton('Add 5 minutes', '+5', 299),
      this.input,
    )
    this.shadowRoot.replaceChildren(style, controls)
    this.shadowRoot.addEventListener('click', this.onClick)
    this.renderDuration()
  }

  disconnectedCallback() {
    this.shadowRoot.removeEventListener('click', this.onClick)
    this.input?.removeEventListener('change', this.onInputChange)
    this.input?.removeEventListener('keydown', this.onInputKeyDown)
    this.input?.removeEventListener('blur', this.onInputBlur)
    this.removeTimerTickSubscription()
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (name === 'data-timer-id' && oldValue !== newValue && this.isConnected) {
      this.updateTimerTickSubscription()
    }
  }

  updateTimerTickSubscription() {
    this.removeTimerTickSubscription()

    const timerId = this.getAttribute('data-timer-id')
    if (!timerId) {
      return
    }

    this.timerElement = document.getElementById(timerId)
    this.timerElement?.addEventListener(presentationTimerTickEventType, this.onTimerTick)
  }

  removeTimerTickSubscription() {
    this.timerElement?.removeEventListener(presentationTimerTickEventType, this.onTimerTick)
    this.timerElement = null
  }

  onInputKeyDown(event) {
    if (event.key !== 'Enter') {
      return
    }

    event.preventDefault()
    this.onInputChange()
    this.input.blur()
    this.renderDuration()
  }

  onInputChange() {
    const duration = this.parseDuration(this.input.value)
    if (duration === null) {
      return
    }

    this.duration = duration
    if (this.timerElement) {
      this.timerElement.ElapsedDuration = 0
    }
    this.updateTimerDuration()
  }

  onInputBlur() {
    this.renderDuration()
  }

  parseDuration(value) {
    const match = /^(\d+)(?:\s*[\:\-\.]\s*(\d+))?$/.exec(value.trim())
    if (!match) {
      return null
    }

    const minutes = Number(match[1])
    const seconds = Number(match[2] ?? 0)
    const duration = minutes * 60 + seconds
    return Number.isSafeInteger(duration) ? duration : null
  }

  onTimerTick(event) {
    if (event.currentTarget !== this.timerElement) {
      return
    }

    const remainingDuration = Number(event.detail?.remainingDuration)
    if (!Number.isFinite(remainingDuration) || remainingDuration < 0) {
      return
    }

    this.duration = Math.ceil(remainingDuration / 1000)
    this.renderDuration()
  }

  getDurationAttribute() {
    if (this.timerElement) {
      const timerDuration = Number(this.timerElement.UtmostDuration)
      if (Number.isFinite(timerDuration) && timerDuration >= 0) {
        return Math.floor(timerDuration / 1000)
      }
    }

    const value = Number(this.getAttribute('duration'))
    return Number.isFinite(value) && value >= 0 ? Math.floor(value) : 60
  }

  createButton(label, text, adjustment) {
    const button = document.createElement('button')
    button.type = 'button'
    button.textContent = text
    button.setAttribute('aria-label', label)
    button.dataset.adjustment = String(adjustment)
    return button
  }

  onClick(event) {
    const button = event.target.closest('button[data-adjustment]')
    if (!button) {
      return
    }

    const adjustment = button.dataset.adjustment
    if (adjustment === 'reset') {
      this.duration = this.defaultDuration
    } else {
      this.duration = Math.max(0, this.duration + Number(adjustment))
      // round up to the nearest minute
      // this.duration = Math.ceil(this.duration / 60) * 60
    }
    this.updateTimerDuration()
    this.renderDuration()
  }

  updateTimerDuration() {
    if (this.timerElement) {
      this.timerElement.UtmostDuration = this.duration * 1000
      this.timerElement.start()
    }
  }

  renderDuration() {
    if (!this.input || this.shadowRoot.activeElement === this.input) {
      return
    }

    const minutes = Math.floor(this.duration / 60)
    const seconds = this.duration % 60
    this.input.value = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
    this.input.setAttribute('aria-valuetext', `${minutes} minutes ${seconds} seconds`)
  }
}

if (!customElements.get('timer-set')) {
  customElements.define('timer-set', TimerSet)
}
