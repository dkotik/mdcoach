<script>
  export let index = 0
  export let active = false
  export let visible = false

  import { onMount } from 'svelte'
  let scaled = false
  let element

  import { tick } from 'svelte'
  const scale = async () => {
    // const parentHeight = element.parentNode.offsetHeight
    const parentHeight = element.parentNode.parentNode.clientHeight
    let ratio = parentHeight / (element.clientHeight + 20)
    if (ratio >= 1) return // already scaled

    let fontSize = 100
    while (fontSize > 50) {
      fontSize -= 10
      element.style.fontSize = fontSize + '%'
      await tick()
      ratio = parentHeight / (element.clientHeight + 20)
      // console.log("finished scaling", element.style.fontSize, element.style.transform)
      if (ratio >= 1) return // already scaled
    }
    element.style.transform = 'scale(' + ratio + ')'
  }

  $: if(element && visible && !scaled) {
    scaled = true
    scale()
  }

  const unlockScaling = (event) => {
    tick().then(() => scaled = false)
    console.log("scaling changed:", scaled)
  }
  onMount(() => {
    window.addEventListener("resize", unlockScaling)
    return () => window.removeEventListener("resize", unlockScaling)
  })

  // import { slideIn, slideOut } from '../navigation/transitions.js'
  // in:slideIn={{ duration: 800 }}
  // out:slideOut={{ duration: 800 }}
</script>

<section
  class:active={active}
>
  <article bind:this={element}>
    {#if visible}
      <slot />
      <!-- <div style="width: 200vw; background-color:purple;">&nbsp;</div> -->
      <!-- <div style="height: 200vh; background-color:purple;">&nbsp;</div> -->
      <p>[{index}]</p>
      <p>
        {#each Array(Math.ceil(Math.random() * 1000)) as index}
          {index} .
        {/each}
      </p>
    {:else}
      TODO: LOADING...
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
  color: white;
  margin: 0 auto;
  display: grid;
  /* grid-auto-flow: row; */
  /* grid-template-rows: max-content; */
  grid-template-columns: 1fr auto 1fr;
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
