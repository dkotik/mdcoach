<script>
  import './Daggers.css'
  import Keys from '../navigation/Keys.svelte'
  import { Dispatch } from '../navigation/broadcast.js'
  import { verticalScrollTo, isVerticalScrollNecessary } from '../navigation/scroll.js'
  import { onMount } from 'svelte'
  export let slideData
  export let active = 1

  // TODO: add the necessity check of vertical scroll
  // setTimeout(() => {

  const onSlideChange = (event) => {
    active = event.slide
    // verticalScrollTo('divider'+event.slide)
    if (isVerticalScrollNecessary('marker'+event.slide)) verticalScrollTo('divider'+event.slide)
  }
  onMount(() => {
    // const velement = document.getElementById("marker11")
    // setInterval(() => {
    //   console.clear()
    //   if (isVerticalScrollNecessary(velement)) console.log("vnec")
    // }, 500)

    verticalScrollTo('divider'+active)
    window.addEventListener('slideChange', onSlideChange)
    return () => window.removeEventListener('slideChange', onSlideChange)
  })
</script>

<Keys
  daggerQuery='body > #app > main > div.notes > section.active > ul > li'
  on:dagger={(event) => Dispatch(active, event.detail.number)}
  on:next={() => {
    if (active < (slideData.slides || []).length) Dispatch(active+1)
  }}
  on:previous={() => {
    if (active > 1) Dispatch(active-1)
  }}
/>

<div class="notes" role="presentation">
{#each slideData.slides as slide, index}
  {@const ID = index + 1}
  <div class="divider" id={'divider'+ID} />
  <aside class="media">
    pictures
  </aside>
  <a
    role="tab"
    tabindex={ID}
    class="marker"
    id={'marker'+ID}
    class:active={active === ID}
    on:mouseup={() => {
      if (isVerticalScrollNecessary('marker'+ID)) verticalScrollTo('divider'+ID)
      Dispatch(ID)
    }}
  >
    {ID}
  </a>
  <section class:active={active === ID}>
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
