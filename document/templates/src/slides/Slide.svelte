<script>
  export let index = 0
  export let active = false
  export let visible = false

  import scaling from './scaling.mjs'
  import { onMount } from 'svelte'
  let element

  // import { slideIn, slideOut } from '../navigation/transitions.js'
  // in:slideIn={{ duration: 800 }}
  // out:slideOut={{ duration: 800 }}
</script>

<section
  class:active={active}
>
  <article bind:this={element} use:scaling={visible}>
    {#if visible}
      <slot />
      <!-- <div style="width: 200vw; background-color:purple;">&nbsp;</div> -->
      <!-- <div style="height: 200vh; background-color:purple;">&nbsp;</div> -->
      <p>[{index}]</p>
      <!-- <p>
        {#each Array(Math.ceil(Math.random() * 1000)) as value, index}
          {index+1} test word.
        {/each}
      </p> -->
    {:else}
      <div class="loading">...</div>
    {/if}
  </article>
</section>
