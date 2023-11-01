import ID from './documentid.js'

const localStorageKey = 'mdcoachCurrentSlide'

export class SlideChangeEvent extends Event {
  presentation = 'unknown'
  slide = 0
  dagger = 0

  constructor(slide=1, dagger=0, presentation=ID) {
    super('slideChange')
    this.presentation = presentation || ID
    this.slide = slide || 1
    this.dagger = dagger || 0
    return this
  }

  get detail() {
    return this.slide
  }

  get anchor() {
    return this.dagger > 0
      ? `#${this.slide}.${this.dagger}`
      : `#${this.slide}`
  }
}

let lastDispatch = new SlideChangeEvent(0, 0)
export function Dispatch(slide, dagger=0) {
  // if (lastDispatch.presentation !== ID) return
  if (lastDispatch.slide === slide && lastDispatch.dagger === dagger) return

  const event = new SlideChangeEvent(slide, dagger)
  lastDispatch = event
  window.dispatchEvent(event)
}

window.addEventListener('slideChange', (event) => {
  window.localStorage.setItem(localStorageKey, JSON.stringify(event))
  console.log("got custom event", event.detail, JSON.stringify(event))

  // console.log("decided to broadcast", parseInt(slide) === event.slide, parseInt(slide), event.slide)
  window.location.hash = event.anchor
})

const handleHashChange = () => {
  console.log("hash change", window.location.hash)
  const [slide, dagger] = window.location.hash.slice(1).split('.', 2)
  // if (event.presentation !== ID) return
  // const [slide, dagger] = window.location.hash.slice(1).split('.', 2)
  if (parseInt(slide) === lastDispatch.slide) return
  Dispatch(parseInt(slide), parseInt(dagger))
}
window.addEventListener('hashchange', handleHashChange)

window.addEventListener('storage', (event) => {
 if (event.storageArea != window.localStorage) return
 if (event.key === localStorageKey) {
   const data = JSON.parse(event.newValue)
   if (data.presentation == ID) Dispatch(
       data.slide,
       data.dagger
   )
 }
})

document.addEventListener('DOMContentLoaded', (event) => {
  const lastSlideEvent = window.localStorage.getItem(localStorageKey)
  if (lastSlideEvent) {
    try {
      const data = JSON.parse(lastSlideEvent)
      if (data.presentation == ID) {
        return Dispatch(
          data.slide,
          data.dagger,
        )
      }
    } catch (e) {
      console.log("saved slide JSON parsing failed:", e)
    }
  }
  handleHashChange()
})
