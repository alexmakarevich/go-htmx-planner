import { create, toBinary, fromBinary } from "@bufbuild/protobuf";
import { CreateUserParamsSchema, UserListSchema } from "../../proto-es/user_pb";

export const listUsers = async () => {
    const response = await fetch("http://localhost:19999/v2/api/list-users", {
      method: "get",
    });
    const parsedBody = fromBinary(UserListSchema, await response.bytes());
    return parsedBody.users;
};

export const createUser = async ({name, password}: {name: string, password: string}) => {
    const createUser = create(CreateUserParamsSchema, { name, password });
    const body = toBinary(CreateUserParamsSchema, createUser);
    await fetch("http://localhost:19999/v2/create-user", {
      method: "post",
      body,
    })
}