// ============== focus.js =====================

let currentSlide = 0
let currentListItem = 0

const focusedClass = "is-focused"
const slides = document.querySelectorAll('main > section')
const finalSlideIndex = slides.length-1

const navigate = (isForward) => {
  slides[currentSlide].classList.remove(focusedClass)
  currentSlide = currentSlide + (
    isForward ? 1 : -1
  )
  slides[currentSlide].classList.add(focusedClass)
}
navigate(true)

window.addEventListener(nextListItemOrSlide, (event) => {
  if(currentSlide >= finalSlideIndex) {
    return
  }
  navigate(true)
})

window.addEventListener(previousSlide, (event) => {
  if (currentSlide === 0) {
    return
  }
  navigate(false)
})
