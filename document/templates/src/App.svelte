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
  // import { onMount } from 'svelte'
  // onMount()

  // {import.meta.env.MODE}
  // {import.meta.env.VITE_SLIDE_DATA}
  let showNotes = false
  let current = 1

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
    <div style="margin: 0 auto;display: flex; place-items: center;height: 80vh;">
      <Loading />
    </div>
  {:then slideData}
    {#if showNotes}
      <Notes active={current} {slideData} />
    {:else}
      <Slides bind:active={current} {slideData} />
    {/if}
  {/await}
</main>

<style>
:global(body) {
  background-color: var(--color-body-background);
}
</style>
