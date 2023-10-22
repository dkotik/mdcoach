<script>
  import './stylesheets/Daggers.css'
  import Slide from './Slide.svelte'
  import Curtain from './Curtain.svelte'
  import Menu from './navigation/Menu.svelte'
  // import Clock from './time/Clock.svelte'
  import './navigation/keys.js'
  // import { onMount } from 'svelte'
  // onMount()

  let slides = []
  let notes = []
  let footnotes = []
  if (Array.isArray(slideData)) {
    for (let i = 0; i < slideData.length; i += 3) {
      slides.push(slideData[i] || "")
      notes.push(slideData[i+1] || "")
      footnotes.push(slideData[i+2] || "")
    }

    if (import.meta.env.DEV) {
      for (let i = 0; i < parseInt(import.meta.env.VITE_SLIDE_MULTIPLIER); i++) {
        slides = slides.concat(slides)
        notes = notes.concat(notes)
        footnotes = footnotes.concat(footnotes)
      }
    }
  }

  // {import.meta.env.MODE}
  // {import.meta.env.VITE_SLIDE_DATA}
</script>

<Menu />

<main>
  <Curtain />
  <!-- <Clock /> -->
  {#each slides as slide, index}
    <Slide index={index+1}>{@html slide}</Slide>
  {/each}
  Add navigation slide limit with Math.ceil(x/3)
  <br />
  Add error list on the notes view
</main>

<style>

</style>
