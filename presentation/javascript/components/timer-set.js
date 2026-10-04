// use as <timer-set duration="300">; duration is in seconds
// emits 'presentation-timer-change' with the selected duration in seconds

class TimerSet extends HTMLElement {
  #timer

  static get observedAttributes() {
    return ['data-timer-id']
  }

  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onInputChange = this.onInputChange.bind(this)
    this.onInputKeyDown = this.onInputKeyDown.bind(this)
    this.onInputBlur = this.onInputBlur.bind(this)
    this.onTimerTick = this.onTimerTick.bind(this)
    this.onTimerChange = this.onTimerChange.bind(this)
    this.defaultDuration = 60
    this.duration = this.defaultDuration
    this.#timer = null
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
      this.createButton('Reset timer', '↻', () => this.ResetTimer()),
      this.createButton('Subtract one minute', '−1', () => this.AdjustTimer(-61)),
      this.createButton('Add 1 minute', '+1', () => this.AdjustTimer(59)),
      this.createButton('Add 5 minutes', '+5', () => this.AdjustTimer(299)),
      this.input,
    )
    this.shadowRoot.replaceChildren(style, controls)
    this.renderDuration()
  }

  disconnectedCallback() {
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

    this.#timer = document.getElementById(timerId)
    this.#timer?.addEventListener(presentationTimerTickEventType, this.onTimerTick)
    this.#timer?.addEventListener('change', this.onTimerChange)
  }

  removeTimerTickSubscription() {
    this.#timer?.removeEventListener(presentationTimerTickEventType, this.onTimerTick)
    this.#timer?.removeEventListener('change', this.onTimerChange)
    this.#timer = null
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
    this.updateTimerDuration(0)
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
    const timer = event.detail
    if (event.currentTarget !== this.#timer || timer !== this.#timer) {
      return
    }

    this.updateDurationFromTimer(timer)
  }

  onTimerChange(event) {
    const timer = event.detail
    if (event.currentTarget !== this.#timer || timer !== this.#timer) {
      return
    }

    this.updateDurationFromTimer(timer)
  }

  updateDurationFromTimer(timer) {
    if (timer !== this.#timer) {
      return
    }

    const remainingDuration = Number(timer.UtmostDuration - timer.ElapsedDuration)
    if (!Number.isFinite(remainingDuration) || remainingDuration < 0) {
      return
    }

    this.duration = Math.ceil(remainingDuration / 1000)
    if (!this.input) {
      return
    }

    const minutes = Math.floor(this.duration / 60)
    const seconds = this.duration % 60
    this.input.value = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
    this.input.setAttribute('aria-valuetext', `${minutes} minutes ${seconds} seconds`)
  }

  getDurationAttribute() {
    if (this.#timer) {
      const timerDuration = Number(this.#timer.UtmostDuration)
      if (Number.isFinite(timerDuration) && timerDuration >= 0) {
        return Math.floor(timerDuration / 1000)
      }
    }

    const value = Number(this.getAttribute('duration'))
    return Number.isFinite(value) && value >= 0 ? Math.floor(value) : 60
  }

  createButton(label, text, onClick) {
    const button = document.createElement('button')
    button.type = 'button'
    button.textContent = text
    button.setAttribute('aria-label', label)
    button.addEventListener('click', onClick)
    return button
  }

  ResetTimer() {
    this.duration = this.defaultDuration
    this.updateTimerDuration()
    this.renderDuration()
  }

  AdjustTimer(adjustment) {
    this.duration = Math.max(0, this.duration + adjustment)
    // round up to the nearest minute
    // this.duration = Math.ceil(this.duration / 60) * 60
    this.updateTimerDuration()
    this.renderDuration()
  }

  updateTimerDuration(elapsed) {
    this.#timer.setDuration(
      elapsed ?? this.#timer.ElapsedDuration,
      this.duration * 1000,
    )
    this.#timer.dispatchChange()
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
