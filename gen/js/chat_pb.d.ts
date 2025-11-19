import * as jspb from 'google-protobuf'



export class ChatMessage extends jspb.Message {
  getUser(): string;
  setUser(value: string): ChatMessage;

  getMessage(): string;
  setMessage(value: string): ChatMessage;

  getTimestamp(): number;
  setTimestamp(value: number): ChatMessage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChatMessage.AsObject;
  static toObject(includeInstance: boolean, msg: ChatMessage): ChatMessage.AsObject;
  static serializeBinaryToWriter(message: ChatMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChatMessage;
  static deserializeBinaryFromReader(message: ChatMessage, reader: jspb.BinaryReader): ChatMessage;
}

export namespace ChatMessage {
  export type AsObject = {
    user: string,
    message: string,
    timestamp: number,
  }
}

export class StreamMessage extends jspb.Message {
  getMessage(): ChatMessage | undefined;
  setMessage(value?: ChatMessage): StreamMessage;
  hasMessage(): boolean;
  clearMessage(): StreamMessage;

  getEventCase(): StreamMessage.EventCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StreamMessage.AsObject;
  static toObject(includeInstance: boolean, msg: StreamMessage): StreamMessage.AsObject;
  static serializeBinaryToWriter(message: StreamMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StreamMessage;
  static deserializeBinaryFromReader(message: StreamMessage, reader: jspb.BinaryReader): StreamMessage;
}

export namespace StreamMessage {
  export type AsObject = {
    message?: ChatMessage.AsObject,
  }

  export enum EventCase { 
    EVENT_NOT_SET = 0,
    MESSAGE = 1,
  }
}

