<script lang="ts">
  import { listEvents } from "../../api/events";
  import { navigate } from "svelte-routing";
</script>

<h2>Events</h2>
<button onclick={() => navigate("/v2/create-event")}>create new</button>
{#await listEvents()}
  <p>loading events...</p>
{:then events}
  <ul>
    {#each events as event}
      <!-- TODO: why could they be undefined?! check proto file -->
      {#if event.event && event.owner}
        <li>
          <h3>{event.event.title}</h3>
          - {event.event.dateTime
            ? new Date(
                Number(event.event.dateTime.seconds * 1000n)
              ).toLocaleString()
            : ""}
          | onwer: {event.owner.name}
        </li>
        <!-- TODO remove verbose logging -->
      {:else}
        {console.error("WTF", event)}
      {/if}
    {/each}
  </ul>
{:catch err}
  <p>could not load events. :(</p>
{/await}
<ul></ul>
