<script>
  import './layout.css'
  import hiddenListItems from '../navigation/list.js'
  import scrollStop from '../actions/scrollStop.js'
  import navigation from '../actions/navigation.js'
  import Slide from './Slide.svelte'
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
  let daggers = null
</script>

<div
  class="slides"
  role="presentation"
  use:navigation={true}
  on:next={(event) => {
    if (daggers?.next()) {
      console.log("ran into a dagger")
      return
    }
    if (active < (slideData.slides || []).length) active += 1
  }}
  on:previous={() => {
    if (active > 1) active -= 1
  }}
  use:scrollStop
  on:scrollStop={(event) => {
    const current = Math.ceil((event.target.scrollLeft + 0.01) / event.target.clientWidth)
    active = current
    daggers = new hiddenListItems(slidesElement)
    console.log("current slide:", current)
    //setTimeout(() => event.target.nextDagger(), 50)
    //event.target.nextDagger()
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
