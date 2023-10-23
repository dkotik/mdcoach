<script>
  let darkMode = window.localStorage.getItem("darkMode") === "true"
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
    Dark
  {:else}
    Light
  {/if}
</button>

<style>

</style>
