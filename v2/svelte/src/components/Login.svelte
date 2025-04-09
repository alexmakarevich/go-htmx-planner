<script lang="ts">
  import { logIn } from "../api/auth";
  import { navigate } from "svelte-routing";
  import { globalToaster } from "./global/toaster.svelte";

  let name = $state("");
  let password = $state("");
</script>

<!-- TODO: handle case where you're already logged in: -->
<!-- "you're already logged in, log out? go back?" *fields disabled* -->

<h2>log in</h2>

<input type="text" bind:value={name} />
<input type="password" bind:value={password} />

<button
  onclick={async () => {
    const { isSuccess, error, badResponse } = await logIn({ name, password });
    if (isSuccess) {
      // TODO: allow navigating to previously attempted pages via query params
      navigate("/v2/");
    } else if (badResponse) {
      globalToaster.add({ type: "failure", message: badResponse.message });
    } else {
      globalToaster.add({ type: "failure", message: error.message });
    }
  }}>log in</button
>
