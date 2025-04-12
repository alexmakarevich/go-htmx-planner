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

<div style="display: flex; flex-direction: row; justify-content: space-between">
  <h2>Events</h2>
  <button onclick={() => navigate("/v2/create-event")}>create new</button>
</div>

{#await events}
  <p>loading events...</p>
{:then events}
  <ul use:scrollToLastPosition>
    {#each events as event (event.event?.id)}
      <!-- TODO: why could they be undefined?! check proto file -->
      {#if event.event && event.owner}
        <li>
          <article class="event-info">
            <div>
              <h3>{event.event.title}</h3>
              <div>
                <i>
                  {event.event.dateTime
                    ? new Date(
                        Number(event.event.dateTime.seconds * 1000n)
                      ).toLocaleString()
                    : ""}
                </i>
              </div>
              owner: {event.owner.name}
            </div>

            <button
              class="secondary"
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
          </article>
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

<style>
  ul {
    margin-top: 1rem;
    padding: 0;
  }
  li {
    list-style: none;
  }

  .event-info {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }
</style>
