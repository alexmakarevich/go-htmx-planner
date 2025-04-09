import { LoginParamsSchema } from "../../proto-es/auth_pb";
import { ErrorResponseSchema } from "../../proto-es/generic_pb";
import { apiReqRes } from "./api-req";

export const logIn = async ({
  name,
  password,
}: {
  name: string;
  password: string;
}) =>
  apiReqRes(
    {
      path: "/login",
      method: "post",
    },
    {
      input: { name, password },
      inputSchema: LoginParamsSchema,
      outputErrorSchema: ErrorResponseSchema,
    },
  );
