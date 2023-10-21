import { writable } from 'svelte/store'

const current = writable(
  parseInt(window.localStorage.getItem("currentSlide")) || 1,
  // () => { // got some subscribers
  //   window.addEventListener('hashchange',() => {})
  //   return () => { // no more subscribers
  //
  //   }
  // }
)

let timeout = null
current.subscribe((value) => {
  clearTimeout(timeout)
  timeout = setTimeout(() => window.localStorage.setItem("currentSlide", value), 1000)
})

export default current
