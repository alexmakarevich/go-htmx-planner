<script lang="ts">
  import { Router, Link, Route, navigate } from "svelte-routing";
  let { children } = $props();
</script>

<nav>
  <Link to="">home</Link>
  <Link to="login">log in</Link>
</nav>
{@render children?.()}
<footer>some footer stuff © 2077</footer>
