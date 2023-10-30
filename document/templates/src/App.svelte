<script>
  import loadSlides from './slides/load.js'
  import Loading from './navigation/Loading.svelte'
  import './stylesheets/Daggers.css'
  import './themes/default.css'
  import Slides from './slides/Slides.svelte'
  import Notes from './notes/Notes.svelte'
  import Curtain from './Curtain.svelte'
  import Menu from './navigation/Menu.svelte'
  // import Clock from './time/Clock.svelte'
  import './navigation/keys.js'
  // import { onMount } from 'svelte'
  // onMount()

  // {import.meta.env.MODE}
  // {import.meta.env.VITE_SLIDE_DATA}
  let showNotes = false
  let current = 1
  $: console.log("top current:", current)

  window.addEventListener(
    'slideChange',
    (event) => current = event.slide
  )
</script>

<Menu
  on:mode={(event) => showNotes = event.detail === 'notes' }
/>
<Curtain />

<main>
  {#await loadSlides('slideData')}
    <Loading />
  {:then slideData}
    {#if showNotes}
      <Notes active={current} {slideData} />
    {:else}
      <Slides bind:active={current} {slideData} />
    {/if}
  {/await}

  <!-- <Clock /> -->
  Add navigation slide limit with Math.ceil(x/3)
  <br />
  Add error list on the notes view
</main>

<style>
:global(body) {
  background-color: var(--color-body-background);
}
</style>
