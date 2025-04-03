<script lang="ts">
  import { DateInput } from "date-picker-svelte";
  import { createEvent } from "../../api/events";
  let title = $state("");
  let dateTime = $state(new Date());
  let visible = $state(false);
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
  max={new Date("2035-12-31")}
>
  <!-- TODO: fix hardcoded max -->

  <button onclick={() => (visible = false)}>close</button>
</DateInput>

<button
  disabled={title === ""}
  onclick={async () => {
    await createEvent({ title, dateTime });
  }}>create event</button
>

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
