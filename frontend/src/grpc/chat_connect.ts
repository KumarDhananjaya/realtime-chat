// frontend/src/grpc/chat_connect.ts
import { StreamMessage } from "./chat_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * The Chat Service definition
 *
 * @generated from service chat.ChatService
 */
export const ChatService = {
  typeName: "chat.ChatService",
  methods: {
    /**
     * A bi-directional stream for sending and receiving messages
     *
     * @generated from rpc chat.ChatService.ChatStream
     */
    chatStream: {
      name: "ChatStream",
      I: StreamMessage,
      O: StreamMessage,
      kind: MethodKind.BiDiStreaming,
    },
  }
} as const;