<script>
  import './layout.css'
  import Slide from './Slide.svelte'
  import Keys from '../navigation/Keys.svelte'
  import { onMount } from 'svelte'
  export let slideData
  export let active = 1

  let isScrolling
  let slidesElement
  $: {
    if (slidesElement) slidesElement.scrollTo((active - 1) * slidesElement.clientWidth, 0)
  }

  // setInterval(() => {
  //   slidesElement.scrollTo(1000, 0)
  // }, 1000)
</script>

<Keys
  daggerQuery='body > #app > main > div.slides > section.active > ul > li'
  on:dagger={(event) => console.log(event.detail.number)}
  on:next={() => {
    if (active < (slideData.slides || []).length) active += 1
  }}
  on:previous={() => {
    if (active > 1) active -= 1
  }}
/>

<div class="slides" role="presentation" on:scroll={(event) => {
  clearTimeout(isScrolling)
  isScrolling = setTimeout(() => {
    const current = Math.ceil((event.target.scrollLeft + 0.01) / event.target.clientWidth)
    active = current
    console.log("current slide:", current)
  }, 100)
}} bind:this={slidesElement}>
{#each slideData.slides as slide, index}
  <!-- <Slide index={index+1}>??{@html slide}</Slide> -->
  <div>{index+1}smdsf ||</div>
{:else}
  TODO: THERE ARE NO SLIDES
{/each}
</div>

<style>

</style>
