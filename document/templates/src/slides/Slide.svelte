<script>
  export let index = 0

  import currentSlide from '../navigation/current.js'
  import { slideIn, slideOut } from '../navigation/transitions.js'
  let isCurrent = false
  $: isCurrent = $currentSlide === index
</script>

{#if isCurrent}
<section
  class:isCurrent={"current"}
  in:slideIn={{ duration: 800 }}
  out:slideOut={{ duration: 800 }}
>
  [{index} = <slot />]
</section>
{/if}

<style>
section {
  width: 12em;
  height: 8em;
  background-color: gray;
}

/* Make each slide a grid layout */
/* article > * {
  position: relative;
  display: grid;
  grid-auto-flow: row;
  grid-template-rows: max-content;
  row-gap: 1rem;
  padding: 2rem;
  font-size: 1.5rem;
} */
</style>
