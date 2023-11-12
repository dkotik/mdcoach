<script>
  import Slide from './Slide.svelte'
  import { keyboardNavigation, wheelNavigation, revealedListItems, scrollStop } from '../controls.mjs'
  import { createEventDispatcher } from 'svelte'
  export let slideData
  export let currentSlide
  export let currentListItem
  const dispatch = createEventDispatcher()
  // let lastActive = 0
  // let isMovingRight = true
  // $: {
  //   isMovingRight = currentSlide > lastActive
  //   lastActive = currentSlide
  //   console.log("moving right:", isMovingRight)
  // }

  let slidesElement
  let scrollTimeout = null
  const scrollToPosition = (x, delay=15, limit=20) => {
    const deltaX = (slidesElement.scrollLeft - x) * 0.8
    if ((deltaX < 5 && deltaX > -5) || limit < 1) {
      slidesElement.scrollTo(x, 0)
      return
    }
    slidesElement.scrollTo(x + deltaX, 0)
    clearTimeout(scrollTimeout)
    scrollTimeout = setTimeout(() => scrollToPosition(x, delay, limit-1), delay)
  }

  $: {
    if (slidesElement) scrollToPosition((currentSlide - 1) * slidesElement.clientWidth)
  }
</script>

<div
  class="slides"
  role="presentation"
  use:keyboardNavigation
  use:wheelNavigation
  use:revealedListItems={[currentSlide, currentListItem]}
  on:previous={() => {
    dispatch("change", {slide: currentSlide-1, listItem: 0})
  }}
  on:nextListItem={(event) => {
    dispatch("change", {slide: event.detail.slide, listItem: event.detail.listItem})
  }}
  on:jump={(event) => {
    dispatch("change", {slide: event.detail, listItem: 0})
  }}
  use:scrollStop
  on:scrollStop={(event) =>
    dispatch("change", {slide: Math.ceil((event.target.scrollLeft + 0.01) / event.target.clientWidth), listItem: 0})
  }
  bind:this={slidesElement}>
{#each slideData.slides as slide, index}
  {@const ID = index+1}
  <Slide
    active={currentSlide === ID}
    visible={currentSlide > ID - 6 && currentSlide < ID + 2}
  >
    {@html slide}
  </Slide>
{:else}
  TODO: THERE ARE NO SLIDES
{/each}
</div>

<style>
.slides {
  font-size: 160%;
  /* background-color: red; */
  width: 100vw;
  height: 100vh;
  display: flex;
  /* place-items: center; */
  overflow-y: hidden;
  overflow-x: scroll;
  /* https://codepen.io/knowler/pen/eYGRwyb */
  /* scroll-snap-type: x mandatory; */
  scrollbar-color: transparent transparent; /* Firefox */
  scrollbar-width: thin;
  /* @supports (overflow-inline: scroll) {
    overflow-inline: scroll;
    scroll-snap-type: inline mandatory;
  } */
}

.slides::-webkit-scrollbar {
  /* Chrome, Safari */
  display: none;
}
</style>
