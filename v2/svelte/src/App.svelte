<script lang="ts">
  import Home from "./Home.svelte";
  import Login from "./Login.svelte";
  import { Router, Link, Route, navigate, useHistory } from "svelte-routing";
  import NotFound from "./NotFound.svelte";
  import LoggedInFrame from "./LoggedInFrame.svelte";

  const queryParams = new URLSearchParams(window.location.search);
  const redirect = queryParams.get("fe-route");
  if (redirect) {
    navigate(redirect, { replace: true });
  }
</script>

<Router basepath="/v2/">
  <Route path="login"><Login /></Route>

  <Route path="/*">
    <Router>
      <Route path="/">
        <!--  -->
        <LoggedInFrame>
          <Router>
            <Route path="/">
              <Home />
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
