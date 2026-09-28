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
        color: var(--color-menu-text, #666);
        display: inline-flex;
        max-width: 100%;
      }

      .controls {
        align-items: center;
        display: flex;
        gap: 0.5rem;
      }

      .group {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
      }

      button {
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-body-subtext, #aaa);
        border-radius: 0.3rem;
        color: inherit;
        cursor: pointer;
        font: inherit;
        line-height: 1.2;
        padding: 0.3rem 0.5rem;
        white-space: nowrap;
      }

      button:hover {
        border-color: var(--color-marker-background, #006eff);
      }

      button:focus-visible {
        outline: 2px solid var(--color-marker-background, #006eff);
        outline-offset: 2px;
      }

      output {
        align-items: center;
        background: var(--color-menu-background, #eee);
        border: 1px solid var(--color-body-subtext, #aaa);
        border-radius: 0.35rem;
        display: flex;
        font-variant-numeric: tabular-nums;
        font-weight: 700;
        gap: 0.2rem;
        padding: 0.45rem 0.6rem;
      }

      .slot {
        min-width: 2ch;
        text-align: center;
      }

      .separator,
      .unit {
        color: var(--color-body-subtext, #aaa);
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

    this.output = document.createElement('output')
    this.output.setAttribute('aria-label', 'Timer duration')

    this.minutes = document.createElement('span')
    this.minutes.className = 'slot'
    this.minutes.setAttribute('aria-label', 'minutes')

    const separator = document.createElement('span')
    separator.className = 'separator'
    separator.textContent = ':'
    separator.setAttribute('aria-hidden', 'true')

    this.seconds = document.createElement('span')
    this.seconds.className = 'slot'
    this.seconds.setAttribute('aria-label', 'seconds')

    // const minuteUnit = document.createElement('span')
    // minuteUnit.className = 'unit'
    // minuteUnit.textContent = 'm'
    // minuteUnit.setAttribute('aria-hidden', 'true')
    // const secondUnit = document.createElement('span')
    // secondUnit.className = 'unit'
    // secondUnit.textContent = 's'
    // secondUnit.setAttribute('aria-hidden', 'true')
    // this.output.append(this.minutes, minuteUnit, separator, this.seconds, secondUnit)
    this.output.append(this.minutes, separator, this.seconds)

    const addButtons = document.createElement('div')
    addButtons.className = 'group'
    addButtons.append(
      this.createButton('Add 1 minute', '+1 min', 60),
      this.createButton('Add 5 minutes', '+5 min', 300),
    )

    const otherButtons = document.createElement('div')
    otherButtons.className = 'group'
    otherButtons.append(
      this.createButton('Reset timer', 'Reset', 'reset'),
      this.createButton('Subtract 30 seconds', '−30 sec', -30),
    )

    const controls = document.createElement('div')
    controls.className = 'controls'
    controls.append(addButtons, this.output, otherButtons)
    this.shadowRoot.replaceChildren(style, controls)
    this.shadowRoot.addEventListener('click', this.onClick)
    this.renderDuration()
  }

  disconnectedCallback() {
    this.shadowRoot.removeEventListener('click', this.onClick)
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
      const timerDuration = Number(this.timerElement.duration)
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
    }
    this.renderDuration()
    if (this.timerElement) {
      this.timerElement.duration = this.duration * 1000
    }
  }

  renderDuration() {
    if (!this.minutes || !this.seconds) {
      return
    }

    const minutes = Math.floor(this.duration / 60)
    const seconds = this.duration % 60
    this.minutes.textContent = String(minutes).padStart(2, '0')
    this.seconds.textContent = String(seconds).padStart(2, '0')
    this.output.setAttribute('aria-label', `${minutes} minutes ${seconds} seconds`)
  }
}

if (!customElements.get('timer-set')) {
  customElements.define('timer-set', TimerSet)
}
