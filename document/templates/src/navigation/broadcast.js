import { get } from 'svelte/store'
import current from './current.js'
import ID from './documentid.js'

const localStorageKey = 'mdcoachCurrentSlide'

const respondToSlideBroadcast = (message) => {
  // if (typeof message !== 'string') return
  const [broadcastID, slideNumber] = message.split(":", 2)
  if (parseInt(broadcastID) !== ID) return
  // TODO: check for top boundary
  current.set(parseInt(slideNumber))
}

window.addEventListener('storage', (event) => {
 if (event.storageArea != window.localStorage) return
 if (event.key === localStorageKey) {
   respondToSlideBroadcast(event.newValue || '')
 }
})

// rewind to the last used slide
respondToSlideBroadcast(window.localStorage.getItem(localStorageKey) || '')

current.subscribe((value) => {
  window.localStorage.setItem(localStorageKey, ID + ":" + value)
  // console.log("broadcast", localStorageKey, ID + ":" + value)
})
