import currentSlide from './current.js'
import daggers from './daggers.js'
import { isCurtainVisible } from '../Curtain.svelte'
import { get } from 'svelte/store'

document.addEventListener('keydown', (event) => {
  if (get(isCurtainVisible)) return
  switch (event.code) {
  // if (e.code == "Escape") {
  //     this.navigator = false;
  //     this.curtain = false;
  //     e.preventDefault();
  //     return false;
  //   } else if (e.code == "Period") {
  //     this.curtain = !this.curtain;
  //     this.navigator = false;
  //   }
  case "ArrowRight":
  case "ArrowDown":
  case "PageDown":
  case "Space":
  case "KeyJ":
    for (const dagger of get(daggers)) {
      if (!dagger.wasRevealed) {
        dagger.classList.add("revealed")
        dagger.wasRevealed = true
        return
      }
    }
    currentSlide.update(value => value+1)
    return
  case "ArrowLeft":
  case "ArrowUp":
  case "PageUp":
  case "Backspace":
  case "KeyK":
    currentSlide.update(value => value > 1 ? value-1 : 1)
    return
  default:
    // console.log(event.code)
  }
})
