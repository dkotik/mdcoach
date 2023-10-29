<script>
  export let daggerQuery = 'section.active > ul > li'
  let daggers = []

  import { Dispatch } from './broadcast.js'
  import { isCurtainVisible } from '../Curtain.svelte'
  import { createEventDispatcher, onMount } from 'svelte'
  import { get } from 'svelte/store'
  const dispatch = createEventDispatcher()

  const handleKeyDownEvent = (event) => {
    if (get(isCurtainVisible)) {
      event.preventDefault()
      return
    }

    switch (event.code) {
      case 'ArrowRight':
      case 'ArrowDown':
      case 'PageDown':
      case 'Space':
      case 'KeyJ':
        event.preventDefault()
        let i = 0
        for (const dagger of daggers) {
          i++
          if (!dagger.wasRevealed) {
            dagger.classList.add('revealed')
            dagger.wasRevealed = true
            dispatch('dagger', {
              number: i,
              element: dagger,
            })
            return
          }
        }
        dispatch('next', event)
        return
      case 'ArrowLeft':
      case 'ArrowUp':
      case 'PageUp':
      case 'Backspace':
      case 'KeyK':
        // console.log("previous")
        event.preventDefault()
        dispatch('previous', event)
        return
      case 'KeyC':
        window.open(window.location.href, '_blank')
        return
    }
  }

  let rememberedDigits = ''
  const handleKeyUpEvent = (event) => {
    switch (event.code) {
      case 'Digit0': rememberedDigits += '0'; return
      case 'Digit1': rememberedDigits += '1'; return
      case 'Digit2': rememberedDigits += '2'; return
      case 'Digit3': rememberedDigits += '3'; return
      case 'Digit4': rememberedDigits += '4'; return
      case 'Digit5': rememberedDigits += '5'; return
      case 'Digit6': rememberedDigits += '6'; return
      case 'Digit7': rememberedDigits += '7'; return
      case 'Digit8': rememberedDigits += '8'; return
      case 'Digit9': rememberedDigits += '9'; return
      case 'KeyG':
        Dispatch(parseInt(rememberedDigits))
      default:
        rememberedDigits = ''
        daggers = document.querySelectorAll(daggerQuery)
        // console.log("found daggers", daggers.length)
        // setInterval(() => {
        // }, 1000)
    }
  }


  onMount(() => {
    daggers = document.querySelectorAll(daggerQuery)
    document.addEventListener('keydown', handleKeyDownEvent)
    document.addEventListener('keyup', handleKeyUpEvent)
    return () => {
      document.removeEventListener('keydown', handleKeyUpEvent)
      document.removeEventListener('keyup', handleKeyUpEvent)
    }
  })
</script>
