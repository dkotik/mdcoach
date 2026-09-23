// ============== slides.js =====================

let currentSlide = 0
let currentListItem = 0

const focusedClass = "is-focused"
const reverseClass = "reverse"
const slides = document.querySelectorAll('main > section')
const finalSlideIndex = slides.length - 1
slides[currentSlide].classList.add(focusedClass)
slides[currentSlide].classList.add(reverseClass)

const getCurrentConcealedListItems = () => {
  return slides[currentSlide].querySelectorAll(":scope > .grid > .content > ul > li:not(.is-revealed)")
}

const navigate = (targetSlide) => {
  const isForward = targetSlide > currentSlide

  let currentSlideClassList = slides[currentSlide].classList
  currentSlideClassList.remove(focusedClass)
  if (isForward) {
    currentSlideClassList.add(reverseClass)
  } else {
    currentSlideClassList.remove(reverseClass)
  }
  currentSlide = targetSlide

  currentSlideClassList = slides[currentSlide].classList
  currentSlideClassList.add(focusedClass)
  if (isForward) {
    currentSlideClassList.remove(reverseClass)
  } else {
    currentSlideClassList.add(reverseClass)
  }
  window.history.replaceState(null, null, '#' + (currentSlide + 1))
  document.title = `${currentSlide+1}/${finalSlideIndex+1}`
}

const navigationCompleteEventType = "slideNavigationFinished"
const dispatchNavigationCompleteEvent = (concealedListItemCount) => {
  window.dispatchEvent(
    new CustomEvent(
      navigationCompleteEventType,
      {
        detail: {
          slideIndex: currentSlide,
          concealedListItemCount: concealedListItemCount
        }
    })
  )
}

window.addEventListener(nextEventType, (event) => {
  if(currentSlide >= finalSlideIndex) {
    return
  }
  let concealedListItems = getCurrentConcealedListItems()
  for (const element of concealedListItems) {
    element.classList.add("is-revealed")
    debounce(dispatchNavigationCompleteEvent, 100)(
      concealedListItems.length-1
    )
    return
  }
  navigate(currentSlide + 1)
  concealedListItems = getCurrentConcealedListItems()
  debounce(dispatchNavigationCompleteEvent, 100)(
    concealedListItems.length
  )
})

window.addEventListener(previousEventType, (event) => {
  if (currentSlide === 0) {
    return
  }
  navigate(currentSlide - 1)
  const concealedListItems = getCurrentConcealedListItems()
  debounce(dispatchNavigationCompleteEvent, 100)(
    concealedListItems.length
  )
})

const getSlideIndexFromLocationHash = () => {
  const rawHash = window.location.hash.substring(1);

  // Parse the string into a numeric value
  let numericValue = parseInt(rawHash, 10);

  // Check if it is a valid number
  if (isNaN(numericValue)) {
    return
  }
  numericValue--
  if (numericValue < 0) {
    return 0
  }
  if (numericValue > finalSlideIndex) {
    return finalSlideIndex
  }
  return numericValue
}

const navigateToWindowHashLocation = () => {
  const slideIndex = getSlideIndexFromLocationHash()
  if (slideIndex === currentSlide) {
    return
  }
  navigate(slideIndex)
}
window.addEventListener(
  "hashchange",
  navigateToWindowHashLocation
)

navigateToWindowHashLocation()
