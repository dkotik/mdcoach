<script>
  export let ID = 'all'
  export let currentSlide = 1
  export let currentListItem = 0

  import { onMount } from 'svelte'
  import { debounce } from './controls.mjs'
  import { createEventDispatcher } from 'svelte'
  const localStorageKey = 'mdcoachSlideAnchor'
  const dispatch = createEventDispatcher()

  const decode = (json) => {
    if (!json) return
    try {
      const event = JSON.parse(json)
      if (!event.presentation) throw new Error("no presentation ID")
      if (event.presentation !== ID) return // not for this document
      if (!event.slide) throw new Error("no slide number")
      dispatch('change', {
        slide: event.slide,
        listItem: event.listItem || 0,
      })
    } catch (e) {
      console.log("could not decode slide broadcast:", json, e)
    }
  }

  const broadcast = debounce((slide, listItem) => {
    window.localStorage.setItem(localStorageKey, JSON.stringify({
      presentation: ID,
      slide: currentSlide,
      listItem: currentListItem
    }))
  }, 200)
  $: broadcast(currentSlide, currentListItem)

  const handleStorageUpdate = (event) => {
    if (event.storageArea != window.localStorage) return
    if (event.key !== localStorageKey) return
    return decode(event.newValue)
  }

  onMount(() => {
    if (!window.location.hash) decode(window.localStorage.getItem(localStorageKey))

    window.addEventListener('storage', handleStorageUpdate)
    return () => {
      window.removeEventListener('storage', handleStorageUpdate)
    }
  })
</script>
