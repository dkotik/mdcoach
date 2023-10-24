<script>
    // https://loading.io/css/
</script>

<div class="spinner">
    <div />
    <div />
    <div />
</div>

<style>
    .spinner {
        display: inline-block;
        position: relative;
        width: 100%;
        height: 8em;
    }
    .spinner div {
        display: inline-block;
        position: absolute;
        /* left: 1em; */
        width: 2em;
        background: var(--color-menu-background);
        border: 1px solid var(--color-menu-text);
        animation: spinner .8s cubic-bezier(0, 0.5, 0.5, 1) infinite;
    }
    .spinner div:nth-child(1) {
        left: 0;
        animation-delay: -0.24s;
    }
    .spinner div:nth-child(2) {
        left: 3em;
        animation-delay: -0.12s;
    }
    .spinner div:nth-child(3) {
        left: 6em;
        animation-delay: 0;
    }
    @keyframes spinner {
        0% {
            top: 1em;
            height: 24em;
        }
        50%,
        100% {
            top: 3em;
            height: 12em;
        }
    }
</style>
