// frontend/src/grpc/chat_pb.ts
import { proto3 } from "@bufbuild/protobuf";

/**
 * Message for sending a chat message
 *
 * @generated from message chat.ChatMessage
 */
export class ChatMessage extends proto3.makeMessageType(
  "chat.ChatMessage",
  () => [
    { no: 1, name: "user", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "message", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "timestamp", kind: "scalar", T: 3 /* ScalarType.INT64 */ },
  ],
) {}

/**
 * A wrapper message for our stream to allow for future event types
 *
 * @generated from message chat.StreamMessage
 */
export class StreamMessage extends proto3.makeMessageType(
  "chat.StreamMessage",
  () => [
    { no: 1, name: "message", kind: "message", T: ChatMessage, oneof: "event" },
  ],
) {}