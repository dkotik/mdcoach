import currentSlide from './current.js'
import { get } from 'svelte/store'

document.addEventListener('keydown', (event) => {
  switch (event.keyCode) {
  case 74: // j, J
    currentSlide.update(value => value+1)
    break
  case 75: // k, K
    currentSlide.update(value => value > 1 ? value-1 : 1)
    break
  }
  // console.log(get(currentSlide), currentSlide)
  console.log(event, currentSlide)
  // console.log(event)
})
