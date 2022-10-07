<script>
    export let TScore
    export let Bmd
</script>

<svelte:options tag="svelte-reports-organs"/>

<label for="tscore">TScore</label>
<input id="tscore" type="number" bind:value={TScore}/>
<label for="bmd">BMD</label>
<input id="bmd" type="number" bind:value={Bmd}/>

<style>

</style>