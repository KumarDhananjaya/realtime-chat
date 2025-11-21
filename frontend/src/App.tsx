import { useState, useEffect } from 'react';
import { JoinScreen } from './components/JoinScreen';
import { ChatScreen } from './components/ChatScreen';

// Imports from the modern @bufbuild libraries
import { createPromiseClient } from "@bufbuild/connect";
import { createGrpcWebTransport } from "@bufbuild/connect-web";

// Correct imports from the generated grpc folder
import { ChatService } from './grpc/chat_connect';
import { ChatMessage, SubscribeRequest } from './grpc/chat_pb';

// Create a transport and a client.
const transport = createGrpcWebTransport({
  baseUrl: 'http://localhost:8080',
});
const client = createPromiseClient(ChatService, transport);

function App() {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [username, setUsername] = useState<string>('');
  const [isJoined, setIsJoined] = useState<boolean>(false);

  const handleJoin = (name: string) => {
    setUsername(name);
    setIsJoined(true);
  };

  // This useEffect hook manages the server stream
  useEffect(() => {
    if (isJoined && username) {
      const controller = new AbortController();

      async function startSubscription() {
        try {
          const request = new SubscribeRequest({ user: username });
          const stream = client.subscribeMessages(request, { signal: controller.signal });

          for await (const res of stream) {
            if (res.event.case === "message" && res.event.value) {
              const msg = res.event.value;
              setMessages(prev => [...prev, msg]);
            }
          }
        } catch (err) {
          // Ignore abort errors as they are expected when leaving/unmounting
          if (err instanceof Error && err.name !== 'AbortError') {
            if ((err as any).name !== 'AbortError') {
              console.error("Stream error:", err);
              alert("Connection to the server was lost.");
              setIsJoined(false);
              setMessages([]);
            }
          }
        }
      }

      startSubscription();

      return () => {
        controller.abort();
      };
    }
  }, [isJoined, username]);

  const handleSendMessage = async (messageText: string) => {
    if (!username) return;

    const chatMsg = new ChatMessage({
      user: username,
      message: messageText,
    });

    try {
      await client.sendMessage(chatMsg);
    } catch (err) {
      console.error("Failed to send message:", err);
      alert("Could not send message.");
    }
  };

  return (
    <main className="bg-gray-100">
      {isJoined ? (
        <ChatScreen username={username} messages={messages.map(m => ({
          user: m.user,
          message: m.message,
          timestamp: Number(m.timestamp)
        }))} onSendMessage={handleSendMessage} />
      ) : (
        <JoinScreen onJoin={handleJoin} />
      )}
    </main>
  );
}

export default App;