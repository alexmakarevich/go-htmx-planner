<script lang="ts">
  import { deleteEvent, listEvents } from "../../api/events";
  import { navigate } from "svelte-routing";
  import { globalToaster } from "../global/toaster.svelte";
  import TrashIcon from "../../icons/TrashIcon.svelte";

  let events = $state(listEvents());
  // events = await listEvents();

  let lastScrollPos = $state(0);

  const refresh = () => {
    events = listEvents();
  };

  const scrollToLastPosition = (node: any) => {
    scrollTo(0, lastScrollPos);
  };
</script>

<h2>Events</h2>
<button onclick={() => navigate("/v2/create-event")}>create new</button>
{#await events}
  <p>loading events...</p>
{:then events}
  <ul use:scrollToLastPosition>
    {#each events as event (event.event?.id)}
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
          <!-- TODO: remove from list after delete -->
          <button
            onclick={async () => {
              const { isSuccess } = await deleteEvent(
                event.event?.id as bigint
              );
              if (isSuccess) {
                refresh();
                globalToaster.add({
                  type: "success",
                  message: "delete succeeded",
                });
                lastScrollPos = window.scrollY;
              } else {
                globalToaster.add({
                  type: "failure",
                  message: "delete failed",
                });
              }
            }}
          >
            <TrashIcon />
          </button>
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
