import {
  fromBinary,
  type Message,
  type MessageInitShape,
  create,
  toBinary,
} from "@bufbuild/protobuf";
import type { GenMessage } from "@bufbuild/protobuf/codegenv1";

// TODO: get from config
const BACKEND_PORT = 19999;
const BASE_URL = "http://localhost:" + BACKEND_PORT + "/v2/api";

// FYI: absolutely insane type overloading is due to the fact that TS can't combine conditional return types
// w/ control flow analysis. Therefore all possible param combos had to be "unrolled" and paired w/ their
// respective return signatures
// always returns, never throws
export function apiReqRes<
  InputShape extends Message,
  OutputShape extends Message,
  OutputErrorShape extends Message,
>(
  reqInfo: Omit<RequestInit, "body"> & { path: string }, // TODO: what else to omit?
  customParams: {
    inputSchema: GenMessage<InputShape>;
    input: MessageInitShape<GenMessage<InputShape>>;
    outputSchema: GenMessage<OutputShape>;
    outputErrorSchema: GenMessage<OutputErrorShape>;
  },
): Promise<
  | {
      isSuccess: true;
      result: OutputShape;
      badResponse?: undefined;
      error?: undefined;
    }
  | {
      isSuccess: false;
      badResponse: OutputErrorShape;
      result?: undefined;
      error?: undefined;
    }
  | {
      isSuccess: false;
      error: Error;
      badResponse?: undefined;
      result?: undefined;
    }
>;
export function apiReqRes<
  OutputShape extends Message,
  OutputErrorShape extends Message,
>(
  reqInfo: Omit<RequestInit, "body"> & { path: string }, // TODO: what else to omit?
  customParams: {
    outputSchema: GenMessage<OutputShape>;
    outputErrorSchema: GenMessage<OutputErrorShape>;
  },
): Promise<
  | {
      isSuccess: true;
      result: OutputShape;
      badResponse?: undefined;
      error?: undefined;
    }
  | {
      isSuccess: false;
      badResponse: OutputErrorShape;
      result?: undefined;
      error?: undefined;
    }
  | {
      isSuccess: false;
      error: Error;
      badResponse?: undefined;
      result?: undefined;
    }
>;
export function apiReqRes<
  InputShape extends Message,
  OutputShape extends Message,
>(
  reqInfo: Omit<RequestInit, "body"> & { path: string }, // TODO: what else to omit?
  customParams: {
    inputSchema: GenMessage<InputShape>;
    input: MessageInitShape<GenMessage<InputShape>>;
    outputSchema: GenMessage<OutputShape>;
  },
): Promise<
  | {
      isSuccess: true;
      result: OutputShape;
      error?: undefined;
    }
  | {
      isSuccess: false;
      error: Error;
      result?: undefined;
    }
>;
export function apiReqRes<InputShape extends Message>(
  reqInfo: Omit<RequestInit, "body"> & { path: string }, // TODO: what else to omit?
  customParams: {
    inputSchema: GenMessage<InputShape>;
    input: MessageInitShape<GenMessage<InputShape>>;
  },
): Promise<
  | {
      isSuccess: true;
      result: void;
      error?: undefined;
    }
  | {
      isSuccess: false;
      error: Error;
      result?: undefined;
    }
>;
export function apiReqRes<
  InputShape extends Message,
  OutputErrorShape extends Message,
>(
  reqInfo: Omit<RequestInit, "body"> & { path: string }, // TODO: what else to omit?
  customParams: {
    inputSchema: GenMessage<InputShape>;
    input: MessageInitShape<GenMessage<InputShape>>;
    outputErrorSchema: GenMessage<OutputErrorShape>;
  },
): Promise<
  | {
      isSuccess: true;
      result: void;
      badResponse?: undefined;
      error?: undefined;
    }
  | {
      isSuccess: false;
      badResponse: OutputErrorShape;
      result?: undefined;
      error?: undefined;
    }
  | {
      isSuccess: false;
      error: Error;
      badResponse?: undefined;
      result?: undefined;
    }
>;
export async function apiReqRes<
  InputShape extends Message,
  OutputShape extends Message,
  OutputErrorShape extends Message,
>(
  reqInfo: Omit<RequestInit, "body"> & { path: string }, // TODO: what else to omit?
  customParams: {
    inputSchema?: GenMessage<InputShape>;
    input?: MessageInitShape<GenMessage<InputShape>>;
    outputSchema?: GenMessage<OutputShape> | undefined;
    outputErrorSchema?: GenMessage<OutputErrorShape>;
  },
) {
  try {
    const { input, inputSchema, outputSchema, outputErrorSchema } =
      customParams;

    let body: Uint8Array | undefined;

    if (input && inputSchema) {
      const message = create(inputSchema, input);
      body = toBinary(inputSchema, message);
    }

    const res = await fetch(BASE_URL + reqInfo.path, {
      ...reqInfo,
      body,
    });

    if (res.ok) {
      if (outputSchema) {
        const parsedBody = fromBinary(outputSchema, await res.bytes());
        // @ts-ignore
        // https://github.com/microsoft/TypeScript/issues/33912
        // https://github.com/microsoft/TypeScript/pull/61136
        return { isSuccess: true, result: parsedBody };
      }
      return { isSuccess: true, result: undefined };
    }

    if (outputErrorSchema) {
      const parsedBody = fromBinary(outputErrorSchema, await res.bytes());
      // @ts-ignore
      // https://github.com/microsoft/TypeScript/issues/33912
      // https://github.com/microsoft/TypeScript/pull/61136
      return { isSuccess: false, badResponse: parsedBody };
    }

    return { isSuccess: false, error: new Error("Unexpected bad response") };
  } catch (error) {
    if (error instanceof Error) {
      return { isSuccess: false, error };
    }
    return { isSuccess: false, error: new Error("Improper error") };
  }
}

// TODO: make this pretty and use it
class BadResponseError {
  status: number;
  statusText: string;

  constructor(res: Response) {
    this.status = res.status;
    this.statusText = res.statusText;
  }
}
