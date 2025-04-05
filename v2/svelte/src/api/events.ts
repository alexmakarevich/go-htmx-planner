import { create, fromBinary, toBinary } from "@bufbuild/protobuf";
import {
  CreateEventParamsSchema,
  ListEventsWithOwnerResponseSchema,
} from "../../proto-es/event_pb";
import { timestampFromDate } from "@bufbuild/protobuf/wkt";

export const listEvents = async () => {
  // TODO: base-url
  const response = await fetch("http://localhost:19999/v2/api/list-events", {
    method: "get",
  });
  const parsedBody = fromBinary(
    ListEventsWithOwnerResponseSchema,
    await response.bytes(),
  );
  return parsedBody.events;
};

export const createEvent = async ({
  title,
  dateTime,
  invitedUserIds,
}: {
  title: string;
  dateTime: Date;
  invitedUserIds?: bigint[];
}) => {
  const createEvent = create(CreateEventParamsSchema, {
    title,
    invitedUserIds,
    dateTime: timestampFromDate(dateTime),
  });
  const body = toBinary(CreateEventParamsSchema, createEvent);
  return await fetch("http://localhost:19999/v2/api/create-event", {
    method: "post",
    body,
  });
};
