// ============== slides.js =====================

let currentSlide = 0
let currentListItem = 0

const focusedClass = "is-focused"
const reverseClass = "reverse"
const slides = document.querySelectorAll('main > section')
const finalSlideIndex = slides.length-1

const navigate = (isForward) => {
  for (const element of slides[currentSlide].querySelectorAll(":scope > .content > ul > li:not(.is-revealed)")) {
    element.classList.add("is-revealed")
    return
  }

  slides[currentSlide].classList.remove(focusedClass)
  if (isForward) {
    slides[currentSlide].classList.add(reverseClass)
    currentSlide++
  } else {
    slides[currentSlide].classList.remove(reverseClass)
    currentSlide--
  }
  slides[currentSlide].classList.add(focusedClass)
  if (isForward) {
    slides[currentSlide].classList.remove(reverseClass)
  } else {
    slides[currentSlide].classList.add(reverseClass)
  }
}
navigate(true)

window.addEventListener(nextEventType, (event) => {
  if(currentSlide >= finalSlideIndex) {
    return
  }
  navigate(true)
})

window.addEventListener(previousEventType, (event) => {
  if (currentSlide === 0) {
    return
  }
  navigate(false)
})
