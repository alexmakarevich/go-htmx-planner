<script lang="ts">
  import { Router, Link, Route, navigate } from "svelte-routing";
  let { children } = $props();
</script>

<header>
  <nav class="main-menu">
    <Link to="">home</Link>
    <Link to="users">users</Link>

    <!-- <Link to="events">events</Link>
    <Link to="my-invites">my invites</Link>
    <Link to="settings">settngs</Link> -->

    <Link to="login">log in</Link>
  </nav>
</header>
<article>
  {@render children?.()}
</article>
<footer>some footer stuff © 2077</footer>
