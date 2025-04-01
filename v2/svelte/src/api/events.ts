import { fromBinary } from "@bufbuild/protobuf";
import { ListEventsWithOwnerResponseSchema  } from "../../proto-es/event_pb";

export const listEvents = async () => {
    // TODO: base-url
    const response = await fetch("http://localhost:19999/v2/api/list-events", {
      method: "get",
    });
    const parsedBody = fromBinary(ListEventsWithOwnerResponseSchema, await response.bytes());
    return parsedBody.events;
};