// use as custom element <presentation-progress> for slide navigation progress

class PresentationProgress extends HTMLElement {
  constructor() {
    super()
    this.attachShadow({ mode: 'open' })
    this.onNavigationComplete = this.onNavigationComplete.bind(this)
  }

  connectedCallback() {
    const style = document.createElement('style')
    style.textContent = `
      :host {
        bottom: 0;
        display: block;
        height: 0.25rem;
        left: 0;
        pointer-events: none;
        position: fixed;
        width: 100vw;
        z-index: 10000;
        background: var(--color-body-subtext, #aaa);
      }


      .fill {
        background: var(--color-marker-background, #006eff);
        height: 100%;
        transition: width 600ms ease;
        width: 0;
      }

      @media (prefers-reduced-motion: reduce) {
        .fill {
          transition: none;
        }
      }
    `

    this.setAttribute('role', 'progressbar')
    this.setAttribute('aria-label', 'Presentation progress')
    this.setAttribute('aria-valuemin', '0')
    this.setAttribute('aria-valuemax', '100')

    this.fill = document.createElement('div')
    this.fill.className = 'fill'
    this.shadowRoot.replaceChildren(style, this.fill)

    this.updateProgress(0)
    window.addEventListener('slideNavigationFinished', this.onNavigationComplete)
  }

  disconnectedCallback() {
    window.removeEventListener('slideNavigationFinished', this.onNavigationComplete)
  }

  onNavigationComplete(event) {
    const { slideIndex, finalSlideIndex: finalSlideIndex } = event.detail
    if (!Number.isInteger(slideIndex) || !Number.isInteger(finalSlideIndex)) {
      return
    }

    const progress = finalSlideIndex > 0
      ? (slideIndex / finalSlideIndex) * 100
      : 100
    // console.log("got event", event.detail)
    this.updateProgress(Math.min(100, Math.max(0, progress)))
  }

  updateProgress(percentage) {
    this.fill.style.width = `${percentage}%`
    this.setAttribute('aria-valuenow', String(Math.round(percentage)))
    // console.log("progress:", percentage)
  }
}

if (!customElements.get('presentation-progress')) {
  customElements.define('presentation-progress', PresentationProgress)
}
