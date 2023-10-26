import './broadcast.js'
import currentSlide from './current.js'
import daggers from './daggers.js'
import { isCurtainVisible } from '../Curtain.svelte'
import { get } from 'svelte/store'

let rememberedDigits = ""
const noteJumpDigit = (event) => {
  switch (event.code) {
  case "Digit0": rememberedDigits += "0"; return
  case "Digit1": rememberedDigits += "1"; return
  case "Digit2": rememberedDigits += "2"; return
  case "Digit3": rememberedDigits += "3"; return
  case "Digit4": rememberedDigits += "4"; return
  case "Digit5": rememberedDigits += "5"; return
  case "Digit6": rememberedDigits += "6"; return
  case "Digit7": rememberedDigits += "7"; return
  case "Digit8": rememberedDigits += "8"; return
  case "Digit9": rememberedDigits += "9"; return
  case "KeyG":
    window.location.hash = "slide" + rememberedDigits
  default:
    rememberedDigits = ""
  }
}

document.addEventListener('keydown', (event) => {
  // console.log(event.code)
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
  case "KeyC":
    window.open(window.location.href, '_blank')
  default:
    noteJumpDigit(event)
    // console.log(event.code)
  }
})
