<script>
  import loadSlides from './slides/load.js'
  import Loading from './navigation/Loading.svelte'
  import './stylesheets/Daggers.css'
  import './themes/default.css'
  import Slides from './slides/Slides.svelte'
  import Notes from './notes/Notes.svelte'
  import Curtain from './Curtain.svelte'
  import LocationHash from './LocationHash.svelte'
  import Menu from './navigation/Menu.svelte'
  // import Clock from './time/Clock.svelte'

  // {import.meta.env.MODE}
  // {import.meta.env.VITE_SLIDE_DATA}
  let showNotes = false
  let currentSlide = 1
  let currentListItem = 0
  const jump = (slides, slide, listItem) => {
    if (slide < 0) {
      slide = 1
    } else if (slide > slides.length) {
      slide = slides.length
    }
    currentSlide = slide
    currentListItem = listItem
  }
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
    <LocationHash
      {currentSlide}
      {currentListItem}
      on:change={(e) => jump(slideData.slides, e.detail.slide, e.detail.listItem)}
    />
    {#if showNotes}
      <Notes
        {slideData}
        {currentSlide}
        {currentListItem}
        on:change={(e) => jump(slideData.slides, e.detail.slide, e.detail.listItem)}
      />
    {:else}
      <Slides
        {slideData}
        {currentSlide}
        {currentListItem}
        on:change={(e) => jump(slideData.slides, e.detail.slide, e.detail.listItem)}
      />
    {/if}
  {/await}
</main>

<style>
:global(body) {
  background-color: var(--color-body-background);
}
</style>
