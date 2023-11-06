<script>
  export let index = 0
  export let active = false
  export let visible = false

  import './markdown/blockquote.css'
  import './markdown/heading.css'
  import './markdown/paragraph.css'
  import './markdown/picture.css'
  import './markdown/aside.css'
  import './markdown/list.css'
  import './markdown/thematicBreak.css'
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
    {#if visible}
      <article bind:this={element} use:scaling={visible}>
      <!-- <figure><picture><source media="(max-width: 128px) and (max-height: 128px)" srcset="https://blog.paavo.me/cities-skylines-2-performance/cs2_character.jpeg"><source srcset="https://blog.paavo.me/cities-skylines-2-performance/cs2_character.jpeg"><img src="https://blog.paavo.me/cities-skylines-2-performance/cs2_character.jpeg" width="800" height="517" alt="not found image" title="test image"></picture><figcaption>test image</figcaption></figure> -->
      <!-- <aside>
???
      </aside> -->
      <slot />
      <!-- <div style="width: 200vw; background-color:purple;">&nbsp;</div> -->
      <!-- <div style="height: 200vh; background-color:purple;">&nbsp;</div> -->
      <!-- <p>[{index}]</p> -->
      <!-- <p>
        {#each Array(Math.ceil(Math.random() * 1000)) as value, index}
          {index+1} test word.
        {/each}
      </p> -->
      </article>
    {:else}
      <article class="loading">...</article>
    {/if}
</section>

<style>
section {
  display: flex;
  /* Fill the scroll container */
  /* overflow: hidden; */
  flex: 0 0 100%;
  /* Scroll to the start of the slide */
  /* scroll-snap-align: start; */
  /* Increment the slide counter */
  /* counter-increment: slide; */
}

section > article {
  /* position: relative; */
  color: white;
  width: 100%;
  /* margin: 0 auto; */
  /* padding: 2em; */
  display: grid;
  /* grid-auto-flow: row; */
  /* grid-template-rows: max-content; */
  grid-template-columns: auto 1fr auto 1fr auto;
}

section > article div.loading {
  /* width: 90vw; */
  /* height: 90vh; */
  /* background-color: var(--color-menu-background); */
}

section > article > :global(*) {
  grid-column: 3;
  /* max-width: 20em; */
}
</style>
