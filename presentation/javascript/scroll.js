const scrollToVisibleSlide = debounce((event) => {
  if (document.body.classList.contains('slides')) {
    if (window.scrollY !== 0) {
      window.scrollTo(0, 0)
    }
    return
  }

  const slideIndex = event.detail?.slideIndex
  if (!Number.isInteger(slideIndex)) {
    return
  }

  const slide = document.getElementById(`slide-${slideIndex + 1}`)
  if (!slide) {
    return
  }

  const bounds = slide.getBoundingClientRect()
  const isVisible = bounds.top >= 0 &&
    bounds.left >= 0 &&
    bounds.bottom <= window.innerHeight &&
    bounds.right <= window.innerWidth
  if (isVisible) {
    return
  }

  slide.style.scrollMargin = '1rem'
  slide.scrollIntoView({ behavior: 'smooth' })
}, 120)

window.addEventListener(navigationCompleteEventType, scrollToVisibleSlide)
