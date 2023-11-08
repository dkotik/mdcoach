<script>
  import { slide } from 'svelte/transition'
  import { toggleFullScreen } from './fullscreen.js'
  import { onMount, createEventDispatcher } from 'svelte'
  import IconCurtain from './icons/curtain.svg.svelte'
  import IconMap from './icons/map.svg.svelte'
  import IconFullScreen from './icons/fullscreen.svelte'
  import DarkLightToggle from './DarkLightToggle.svelte'
  const dispatch = createEventDispatcher()

  let show = false
  let showNotes = window.localStorage.getItem('mode') === 'notes'
  const toggleNotes = () => {
    showNotes = !showNotes
    if (showNotes) {
      dispatch('mode',  'notes')
      window.localStorage.setItem('mode', 'notes')
    } else {
      dispatch('mode',  'slides')
      window.localStorage.removeItem('mode')
    }
  }

  onMount(() => {
    if (showNotes) dispatch('mode',  'notes')
  })
</script>

<section role="toolbar" tabindex="0"
  on:mouseenter={() => show = true}
  on:mouseleave={() => show = false}
>
  {#if show}
    <nav in:slide={{ duration: 800 }} out:slide={{ duration: 100 }}>
      <button on:click={toggleNotes}>
        {#if showNotes}
          <IconMap /> Slides
        {:else}
          <IconCurtain /> Notes
        {/if}
      </button>
      <DarkLightToggle />
      <button on:click={toggleFullScreen}>
        <IconFullScreen /> Fullscreen
      </button>
    </nav>
  {:else}&nbsp;{/if}
</section>

<style>
section {
  z-index: 1;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  min-height: 4vh;
}

nav {
  display: flex;
  /* flex-direction: row-reverse; */
  padding: .6em .65em;
  margin: 0 0 2em 0;
  background-color: var(--color-menu-background);
  gap: 0.6em;
}

:global(button) {
  color: var(--color-menu-text);
  background-color: var(--color-button-background);
  overflow: hidden;
  border-radius: 0.3em;
  flex: 1;
}

:global(svg.icon) {
  max-height: 0.8em;
  max-width: 0.8em;
  margin-bottom: -0.06em;
}
</style>
