<script>
  import { slide } from 'svelte/transition'
  import { toggleFullScreen } from './fullscreen.js'
  let show = false
</script>

<section
  on:mouseenter={() => show = true}
  on:mouseleave={() => show = false}
>
  {#if show}
    <nav in:slide={{ duration: 800 }} out:slide={{ duration: 100 }}>

      <button on:click={toggleFullScreen}>
        Fullscreen
      </button>
    </nav>
  {/if}
</section>

<style>
section {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  min-height: 6em;
}

nav {
  display: flex;
  flex-direction: row-reverse;
  padding: .5em .9em;
  margin: 0 0 2em 0;
  background-color: green;
}

button {
  overflow: hidden;
}
</style>
