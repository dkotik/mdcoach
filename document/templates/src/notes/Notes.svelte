<script>
  import { Dispatch } from '../navigation/broadcast.js'
  import { verticalScrollTo } from '../navigation/scroll.js'
  import { onMount } from 'svelte'
  export let slideData
  export let active = 1

  // TODO: add the necessity check of vertical scroll
  // setTimeout(() => {
  // setInterval(() => {
  //   if (isVerticalScrollNecessary("slide11")) console.log("visible")
  //   console.clear()
  // }, 500)

  const onSlideChange = (event) => {
    verticalScrollTo('divider'+event.slide)
  }
  onMount(() => {
    verticalScrollTo('divider'+active)
    window.addEventListener('slideChange', onSlideChange)
    return () => window.removeEventListener('slideChange', onSlideChange)
  })
</script>

on change anchor
<div class="notes">
{#each slideData.slides as slide, index}
  {@const ID = index + 1}
  <div class="divider" id={'divider'+ID} />
  <aside class="media">
    pictures
  </aside>
  <a
    class="marker"
    class:active={active === ID}
    on:mouseup={() => {
      Dispatch(ID)
      verticalScrollTo('divider'+ID)
    }}
  >
    {ID}
  </a>
  <section>
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
  grid-column-gap: 1em;
}

.divider {
  margin-top: 0.5em;
  height: 0.5em;
  /* background-color: red; */
  grid-column-start: 1;
  grid-column-end: -1;
}

a.marker {
  color: var(--color-body-background);
  background-color: var(--color-menu-background);
  border-bottom-left-radius: 0.9em;
  border-bottom-right-radius: 0.9em;
  cursor: pointer;
}

a.marker.active {
  background-color: var(--color-marker-background);
}
</style>
