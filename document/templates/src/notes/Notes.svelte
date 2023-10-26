<script>
  import { scrollToID, isVerticalScrollNecessary } from '../navigation/scroll.js'
  export let slideData

  // setTimeout(() => {
  // setInterval(() => {
  //   if (isVerticalScrollNecessary("slide11")) console.log("visible")
  //   console.clear()
  // }, 500)
</script>

on change anchor
<div class="notes">
{#each slideData.slides as slide, index}
  {@const ID = 'slide'+(index+1)}
  <aside class="media">
    pictures
  </aside>
  <a class="marker" href={"#"+ID} id={ID}>
    {index+1}
  </a>
  <section
    on:click={() => scrollToID('slide'+(index+2))}
  >
    {@html slide}
  </section>
{:else}
  TODO: THERE ARE NO SLIDES
{/each}
</div>

<style>
.notes {
  max-width: 80em;
  display: grid;
  grid-template-columns: 1fr 2em 5fr;
  grid-gap: 1em;
}

a.marker {
  color: var(--color-body-background);
  background-color: var(--color-menu-background);
  border-bottom-left-radius: 0.9em;
  border-bottom-right-radius: 0.9em;
}

a.marker:target {
  background-color: var(--color-marker-background);
}
</style>
