<script>
  import { slide } from 'svelte/transition'
  import { toggleFullScreen } from './fullscreen.js'
  import DarkLightToggle from '../themes/DarkLightToggle.svelte'
  let show = false
</script>

<section role="toolbar" tabindex="0"
  on:mouseenter={() => show = true}
  on:mouseleave={() => show = false}
>
  {#if show}
    <nav in:slide={{ duration: 800 }} out:slide={{ duration: 100 }}>
      <DarkLightToggle />
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
  background-color: var(--color-menu-background);
  gap: 1em;
}

button {
  overflow: hidden;
}
</style>
