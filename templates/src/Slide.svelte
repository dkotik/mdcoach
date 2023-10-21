<script>
  export let index = 0

  import currentSlide from './navigation/current.js'
  import { slideIn, slideOut } from './navigation/transitions.js'
</script>

{#if $currentSlide === index}
<section
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
</style>
