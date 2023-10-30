import { tick } from 'svelte'

export default class hiddenListItems {
  remaining = []

  constructor(element, query='section.active > ul > li:not(.revealed)') {
    // this.remaining = Array.from(element.querySelectorAll(query))
    tick().then(() => this.remaining = Array.from(element.querySelectorAll(query)))
  }

  next() {
    if (this.remaining.length === 0) return null
    const nextListItem = this.remaining[0]
    nextListItem.classList.add('revealed')
    this.remaining = this.remaining.slice(1)
    return nextListItem
  }
}

// export default function(element) {
//   const hiddenListItems = Array.from(element.querySelectorAll(
//     'section.active > ul > li:not(.revealed)'))
//   if (hiddenListItems.length === 0) return null
//   const nextListItem = hiddenListItems[0]
//   nextListItem.classList.add('revealed')
//   // element.hiddenListItems = element.hiddenListItems.slice(1)
//   return nextListItem
// }
