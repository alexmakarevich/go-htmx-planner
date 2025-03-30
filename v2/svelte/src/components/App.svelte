<script lang="ts">
  import Home from "./Home.svelte";
  import { Router, Link, Route, navigate, useHistory } from "svelte-routing";
  import LoggedInFrame from "./LoggedInFrame.svelte";
  import Login from "./Login.svelte";
  import NotFound from "./NotFound.svelte";
  import Users from "./users/Users.svelte";

  const queryParams = new URLSearchParams(window.location.search);
  const redirect = queryParams.get("fe-route");
  if (redirect) {
    navigate(redirect, { replace: true });
  }
</script>

<Router basepath="/v2/">
  <!-- FYI: this router mess is due to the way svelte-routing handles fallbacks and generic routes w/ conditional HTML elements -->
  <Route path="login"><Login /></Route>

  <Route path="/*">
    <Router>
      <Route path="/*">
        <!--  -->
        <LoggedInFrame>
          <Router>
            <Route path="/">
              <Home />
            </Route>
            <Route path="/users">
              <Users />
            </Route>
            <Route path="users">
              <Users />
            </Route>
          </Router>
        </LoggedInFrame>
        <!--  -->
      </Route>
      <Route>
        <NotFound />
      </Route>
    </Router>
  </Route>
</Router>

<style>
</style>
