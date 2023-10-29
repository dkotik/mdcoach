<script>
  import Slide from './Slide.svelte'
  import { onMount } from 'svelte'
  export let slideData
</script>

{#each slideData.slides as slide, index}
  <Slide index={index+1}>??{@html slide}</Slide>
{:else}
  TODO: THERE ARE NO SLIDES
{/each}

<style>

</style>
