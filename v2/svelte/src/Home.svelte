<script lang="ts">
  // import { CreateUserParams, UserList } from "../proto-ts/user";
  import { CreateUserParamsSchema, UserListSchema } from "../proto-es/user_pb";
  import { create, toBinary, fromBinary } from "@bufbuild/protobuf";

  const listUsers = async () => {
    const response = await fetch("http://localhost:19999/v2/api/list-users", {
      method: "get",
    });
    const parsedBody = fromBinary(UserListSchema, await response.bytes());
    return parsedBody.users;
  };
</script>

<h1>HOME</h1>
<h2>Users</h2>
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
