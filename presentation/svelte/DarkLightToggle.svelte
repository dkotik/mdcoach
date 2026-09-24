<script>
  import IconDarkMode from './icons/darkmode.svg.svelte'
  let darkMode = window.localStorage.getItem("darkMode") === "true"
  if (!darkMode) document.body.classList.remove("dark")
</script>

<button
  on:click={()=>{
    darkMode = !darkMode
    if (darkMode) {
      window.localStorage.setItem("darkMode", "true")
      document.body.classList.add("dark")
    } else {
      window.localStorage.removeItem("darkMode")
      document.body.classList.remove("dark")
    }
  }}
>
  {#if darkMode}
    <IconDarkMode /> Dark
  {:else}
    <IconDarkMode /> Light
  {/if}
</button>

<style>

</style>
