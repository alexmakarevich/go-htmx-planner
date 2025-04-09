<script lang="ts">
  import { Router, Link, Route, navigate } from "svelte-routing";
  import { logOut } from "../api/auth";
  import { globalToaster } from "./global/toaster.svelte";
  let { children } = $props();
</script>

<header>
  <nav class="main-menu">
    <Link to="">home</Link>
    <Link to="users">users</Link>
    <Link to="events">events</Link>
    <!-- <Link to="my-invites">my invites</Link> -->
    <!-- <Link to="settings">settngs</Link> -->

    <Link
      to="login"
      onclick={async () => {
        const { isSuccess } = await logOut();
        if (isSuccess) {
          globalToaster.add({ message: "logged out" });
        } else {
          globalToaster.add({ message: "failed to log out", type: "failure" });
        }
      }}>log out</Link
    >
  </nav>
</header>
<article>
  {@render children?.()}
</article>
<footer>some footer stuff © 2077</footer>
