<script>
  // import { CreateUserParams, UserList } from "../proto-ts/user";
  import { CreateUserParamsSchema, UserListSchema } from "../proto-es/user_pb";
  import { create, toBinary, fromBinary } from "@bufbuild/protobuf";

  const listUsers = async () => {
    const response = await fetch("http://localhost:19999/v2/list-users", {
      method: "get",
    });
    const parsedBody = fromBinary(UserListSchema, await response.bytes());
    return parsedBody.users;
  };

  let name = $state("");
  let password = $state("");
</script>

<input type="text" bind:value={name} />
<input type="password" bind:value={password} />

<button
  onclick={async () => {
    const createUser = create(CreateUserParamsSchema, { name, password });
    const body = toBinary(CreateUserParamsSchema, createUser);
    await fetch("http://localhost:19999/v2/create-user", {
      method: "post",
      body,
    });
  }}>SEND IT</button
>

<h1>USERS</h1>
{#await listUsers()}
  <p>loading users...</p>
{:then users}
  <ul>
    {#each users as user}
      <li>{user.name}</li>
    {/each}
  </ul>
{:catch err}
  <p>could not load users. :(</p>
{/await}
<ul></ul>
