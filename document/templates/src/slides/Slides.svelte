<script>
  import './layout.css'
  import Slide from './Slide.svelte'
  import { keyboardNavigation, wheelNavigation, revealedListItems, scrollStop } from '../controls.mjs'
  import { onMount } from 'svelte'
  export let slideData
  export let active = 1
  // let lastActive = 0
  // let isMovingRight = true
  // $: {
  //   isMovingRight = active > lastActive
  //   lastActive = active
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
    if (slidesElement) scrollToPosition((active - 1) * slidesElement.clientWidth)
  }
</script>

<div
  class="slides"
  role="presentation"
  use:keyboardNavigation
  use:wheelNavigation
  use:revealedListItems
  on:previous={() => {
    if (active > 1) active -= 1
  }}
  on:next={(event) => {
    if (event.defaultPrevented) return
    if (active < (slideData.slides || []).length) active += 1
  }}
  on:jump={(event) => {
    if (event.detail < (slideData.slides || []).length) active = event.detail
  }}
  use:scrollStop
  on:scrollStop={(event) => {
    const current = Math.ceil((event.target.scrollLeft + 0.01) / event.target.clientWidth)
    active = current
  }}
  bind:this={slidesElement}>
{#each slideData.slides as slide, index}
  {@const ID = index+1}
  <Slide
    index={ID}
    active={active === ID}
    visible={active > ID - 6 && active < ID + 2}
  >
    {@html slide}
  </Slide>
{:else}
  TODO: THERE ARE NO SLIDES
{/each}
</div>

<style>

</style>
