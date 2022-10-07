<script>
import Organ from './Organ.svelte'

export let Organs
</script>
{#each Organs as {Direction, Site, Bmd, TScore}}
<h4>{Direction +" " +Site}</h4>
   <Organ bind:Organs />
{/each}
<style>

</style>