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
      <p>
        {#each Array(Math.ceil(Math.random() * 1000)) as value, index}
          {index+1} test word.
        {/each}
      </p>
    {:else}
      <div class="loading">...</div>
    {/if}
  </article>
</section>

<style>
section {
  overflow: hidden;
  /* display: flex;
  place-items: center; */
}

article {
  /* position: relative; */
  color: white;
  margin: 0 auto;
  padding: 2em;
  display: grid;
  /* grid-auto-flow: row; */
  /* grid-template-rows: max-content; */
  grid-template-columns: 1fr auto 1fr;
}

article div.loading {
  /* width: 90vw; */
  /* height: 90vh; */
  /* background-color: var(--color-menu-background); */
}

article > * {
  grid-column: 2;
  /* max-width: 20em; */
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
