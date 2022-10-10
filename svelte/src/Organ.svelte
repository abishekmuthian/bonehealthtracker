<script>
    export let TScore
    export let Bmd
</script>

<svelte:options tag="svelte-reports-organs"/>

<section>
    <label for="tscore">TScore:</label>
    <input id="tscore" type="number" step="any" min=-5 max=5  bind:value={TScore}/>
    <label for="bmd">BMD:</label>
    <input id="bmd" type="number" step="any" min=0.001 max=2  bind:value={Bmd}/>
</section>

<style>
</style>