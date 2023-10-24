<script>
  import loadSlides from './slides/load.js'
  import Loading from './navigation/Loading.svelte'
  import './stylesheets/Daggers.css'
  import './themes/default.css'
  import Slide from './Slide.svelte'
  import Curtain from './Curtain.svelte'
  import Menu from './navigation/Menu.svelte'
  // import Clock from './time/Clock.svelte'
  import './navigation/keys.js'
  // import { onMount } from 'svelte'
  // onMount()

  // {import.meta.env.MODE}
  // {import.meta.env.VITE_SLIDE_DATA}
</script>

<Menu />
<Curtain />

<main>
  {#await loadSlides('slideData')}
    <Loading />
  {:then slideData}
    {#each slideData.slides as slide, index}
      <Slide index={index+1}>??{@html slide}</Slide>
    {:else}
      TODO: THERE ARE NO SLIDES
    {/each}
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
