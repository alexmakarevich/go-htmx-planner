<script lang="ts">
  import { LoginParamsSchema } from "../../proto-es/auth_pb";

  import { create, toBinary } from "@bufbuild/protobuf";

  let name = $state("");
  let password = $state("");
</script>

<h2>log in</h2>

<input type="text" bind:value={name} />
<input type="password" bind:value={password} />

<button
  onclick={async () => {
    const loginParams = create(LoginParamsSchema, { name, password });
    const body = toBinary(LoginParamsSchema, loginParams);
    await fetch("http://localhost:19999/v2/api/login", {
      method: "post",
      body,
    });
  }}>log in</button
>
