<script context="module">
  import { writable } from 'svelte/store'
  export let isCurtainVisible = writable(false)
  document.addEventListener('keydown', (event) => {
    switch (event.code) {
    case "Escape":
      isCurtainVisible.set(false)
      return
    case "Period":
      isCurtainVisible.update(value => !value)
      return
    }
  })
</script>

{#if $isCurtainVisible}
  <aside />
{/if}

<style>
aside {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: black;
  z-index: 9999998;
}
</style>
