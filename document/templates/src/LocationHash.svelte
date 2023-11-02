<script>
  export let currentSlide = 1
  export let currentListItem = 0

  import { onMount } from 'svelte'
  import { debounce } from './controls.mjs'
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  const updateHash = debounce((slide, listItem) => {
    if (currentListItem === 0) {
      window.location.hash = slide
      return
    }
    window.location.hash = slide + '.' + listItem
  }, 300)
  $: updateHash(currentSlide, currentListItem)

  const handleHashChange = (event) => {
    const [slide, listItem] = window.location.hash.slice(1).split('.', 2)
    dispatch('change', {
      slide: parseInt(slide) || 1,
      listItem: parseInt(listItem) || 0
      // listItem: 0
    })
  }

  onMount(() => {
    handleHashChange(window.location.hash)
    window.addEventListener('hashchange', handleHashChange)
    return () => {
      window.removeEventListener('hashchange', handleHashChange)
    }
  })
</script>
