<script lang="ts">
  import Home from "./Home.svelte";
  import Login from "./Login.svelte";
  import { Router, Link, Route, navigate } from "svelte-routing";
  import NotFound from "./NotFound.svelte";

  let url = $state("");
  let queryParams = new URLSearchParams(window.location.search);
  let redirect = queryParams.get("redirect");
  if (redirect) {
    navigate(redirect, { replace: true });
  }
</script>

<main>
  <Router {url} basepath="/v2">
    <nav>
      <Link to="">home</Link>
      <Link to="login">log in</Link>
    </nav>
    <div>
      <Route path="/"><Home /></Route>
      <Route path="login"><Login /></Route>
      <Route>
        <NotFound />
      </Route>
    </div>
  </Router>
</main>

<style>
</style>
