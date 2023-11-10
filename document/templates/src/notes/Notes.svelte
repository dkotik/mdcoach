<script>
  export let slideData
  export let currentSlide
  export let currentListItem

  import './markdown/list.css'
  import './markdown/heading.css'
  import './markdown/blockquote.css'
  import './markdown/picture.css'
  import { keyboardNavigation, revealedListItems } from '../controls.mjs'
  import { scrollToPosition, isVerticalScrollNecessary } from './scroll.js'
  import { tick, createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()
  // let scrollTarget
  // $: scrollToPosition(0, scrollTarget, 60)
  const updateScrollTarget = (slide) => tick().then(() => {
    if (!isVerticalScrollNecessary('marker'+slide)) return
    const element = document.getElementById('divider'+slide)
    if (!element) throw new Error("divider not found: "+slide)
    // scrollTarget = element.offsetTop
    scrollToPosition(0, element.offsetTop, 60)
  })
  $: updateScrollTarget(currentSlide)
</script>

<div
  class="notes"
  role="presentation"
  use:keyboardNavigation
  use:revealedListItems={[currentSlide, currentListItem]}
  on:previous={() => {
    dispatch("change", {slide: currentSlide-1, listItem: 0})
  }}
  on:nextListItem={(event) => {
    if (isVerticalScrollNecessary('marker'+event.detail.slide)) {
      updateScrollTarget(event.detail.slide)
    }
    dispatch("change", {slide: event.detail.slide, listItem: event.detail.listItem})
  }}
  on:jump={(event) => {
    dispatch("change", {slide: event.detail, listItem: 0})
  }}
>
  <h1>{document.title || '...'}</h1>
  {#each slideData.slides as slide, index}
    {@const ID = index + 1}
    <div class="divider" id={'divider'+ID} />
    <aside class="media"></aside>
    <a
      role="tab"
      tabindex={ID}
      class="marker"
      id={'marker'+ID}
      class:active={currentSlide === ID}
      on:mouseup={() => dispatch("change", {slide: ID, listItem: 0})}
    >
      {ID}
    </a>
    <section class:active={currentSlide === ID}>
      <article>{@html slide}</article>
      {#if slideData.notes[index]}
        <hr class="notes" />
        {@html slideData.notes[index]}
      {/if}
      {#if slideData.footnotes[index]}
        <hr class="footnotes" />
        {@html slideData.footnotes[index]}
      {/if}
    </section>
  {:else}
    TODO: THERE ARE NO SLIDES
  {/each}
</div>

<style>
.notes {
  max-width: 60em;
  margin: 0 auto;
  padding-bottom: 2em;
  display: grid;
  grid-template-columns: 1fr 2em 5fr 1fr;
  grid-column-gap: 1em;
}

section {
  position: relative;
}

hr {
  border: 0;
  height: 0;
  margin: 1em auto;
  max-width: 12em;
  border-top: 1px solid var(--color-menu-background);
}

hr.notes {
  border-top: 2px dashed var(--color-menu-background);
  margin: 0;
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
