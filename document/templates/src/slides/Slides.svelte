<script>
  import './layout.css'
  import hiddenListItems from '../navigation/list.js'
  import scrollStop from '../actions/scrollStop.js'
  import navigation from '../actions/navigation.js'
  import Slide from './Slide.svelte'
  import { onMount } from 'svelte'
  export let slideData
  export let active = 1

  let slidesElement
  $: {
    if (slidesElement) slidesElement.scrollTo((active - 1) * slidesElement.clientWidth, 0)
  }
  let daggers = []
</script>

<div
  class="slides"
  role="presentation"
  use:navigation
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
  <!-- <Slide index={index+1}>??{@html slide}</Slide> -->
  <section
    class:active={active===index+1}
    class:visible={active > index - 5 && active < index + 3}
  >{index+1}
    {active===index+1}
  ||{@html slide}</section>
{:else}
  TODO: THERE ARE NO SLIDES
{/each}
</div>

<style>

</style>
