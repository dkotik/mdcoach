import { writable, get } from 'svelte/store'
import current from './current.js'

const daggers = writable([])

// try {
//   this.daggers = slide.$el.querySelectorAll(":scope > ul > li");
// } catch (e) {
//   this.daggers = [];
// }

let timeout = null
current.subscribe((value) => {
  if (timeout) {
    clearTimeout(timeout)
    timeout = null
  } else {
    daggers.set([])
  }
  timeout = setTimeout(() => {
    const all = document.querySelectorAll("main > section.isCurrent > ul > li")
    // console.log("daggers", all)
    daggers.set(all)
  }, 1200) // slower than the transition animation?
})

// TODO: export nextDagger function!
// TODO: add upcoming class to dagger
export default daggers
