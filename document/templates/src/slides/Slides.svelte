<script>
  import './layout.css'
  import scrollStop from '../actions/scrollStop.js'
  import Slide from './Slide.svelte'
  import { keyboardNavigation, wheelNavigation, revealedListItems } from '../controls.mjs'
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
  $: {
    if (slidesElement) slidesElement.scrollTo((active - 1) * slidesElement.clientWidth, 0)
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
