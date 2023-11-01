<script>
  import './Daggers.css'
  import { keyboardNavigation, revealedListItems } from '../controls.mjs'
  import { Dispatch } from '../navigation/broadcast.js'
  import { verticalScrollTo, isVerticalScrollNecessary } from './scroll.js'
  import { onMount } from 'svelte'
  export let slideData
  export let active = 1

  const onSlideChange = (event) => {
    active = event.slide
    if (isVerticalScrollNecessary('marker'+event.slide)) verticalScrollTo('divider'+event.slide)
  }
  onMount(() => {
    verticalScrollTo('divider'+active)
    window.addEventListener('slideChange', onSlideChange)
    return () => window.removeEventListener('slideChange', onSlideChange)
  })
</script>

Add error list on the notes view

<div
  class="notes"
  role="presentation"
  use:keyboardNavigation
  use:revealedListItems
  on:previous={(event) => {
    if (active > 1) Dispatch(active-1)
  }}
  on:next={(event) => {
    if (isVerticalScrollNecessary('marker'+active)) {
      verticalScrollTo('divider'+active)
      return
    }
    if (event.defaultPrevented) return
    if (active < (slideData.slides || []).length) Dispatch(active+1)
  }}
  on:jump={(event) => console.log("jump to:", event.detail)}
>
  <h1>{document.title || '...'}</h1>
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
      <article>{@html slide}</article>
    </section>
  {:else}
    TODO: THERE ARE NO SLIDES
  {/each}
</div>

<style>
.notes {
  max-width: 60em;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 2em 5fr;
  grid-column-gap: 1em;
}

h1 {
  grid-column: 2 / 4;
  margin: 4vh 0 0 0;
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
  /* border-bottom-left-radius: 0.9em; */
  /* border-bottom-right-radius: 0.9em; */
  border-radius: 0.4em;
  cursor: pointer;
  text-align: center;
  font-weight: bold;
}

a.marker.active {
  background-color: var(--color-marker-background);
}
</style>
