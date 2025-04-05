<script lang="ts">
  import { DateInput } from "date-picker-svelte";
  import { globalToaster } from "../global/toaster.svelte";
  import { createEvent } from "../../api/events";
  import { navigate } from "svelte-routing";
  let title = $state("");
  let dateTime = $state(new Date());
  let visible = $state(false);

  const handleCreate = async () => {
    try {
      const res = await createEvent({ title, dateTime });
      if (res.ok) {
        console.log("OK");
        globalToaster.add({
          type: "success",
          message: "created event",
        });
        navigate("/v2/events");
      } else {
        globalToaster.add({
          type: "failure",
          message: "could not create event",
        });
      }
    } catch (err) {
      globalToaster.add({
        type: "failure",
        message: "could not create event - network error",
      });
    }
  };
</script>

<h1>Create Event</h1>

<input type="text" bind:value={title} />
<!-- TODO: maybe just use the native form input after all -->
<!-- TODO: proper TIME selector too -->
<DateInput
  bind:value={dateTime}
  bind:visible
  timePrecision={"minute"}
  min={new Date()}
  max={new Date(Date.now() + 315576000000)}
>
  <!-- TODO: fix hardcoded max: +10 years in ms -->

  <button onclick={() => (visible = false)}>close</button>
</DateInput>

<button disabled={title === ""} onclick={handleCreate}>create event</button>

<style>
  :root {
    --date-picker-foreground: var(--pico-color);
    --date-picker-background: var(--pico-form-element-background-color);
    --date-picker-highlight-border: var(
      --pico-form-element-active-border-color
    );
    --date-picker-highlight-shadow: default;
    --date-picker-selected-color: var(--pico-primary);
    --date-picker-selected-background: var(
      --pico-form-element-selected-background-color
    );
    --date-input-width: auto;
  }
</style>
