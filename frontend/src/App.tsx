import { useState, useEffect } from 'react';
import { JoinScreen } from './components/JoinScreen';
import { ChatScreen } from './components/ChatScreen';

// Imports from the modern @bufbuild libraries
import { createPromiseClient } from "@bufbuild/connect";
import { createGrpcWebTransport } from "@bufbuild/connect-web";

// FIX: All gRPC-related imports now correctly point to the '/src/gen/' folder
import { ChatService } from './gen/chat_connect';
import { ChatMessage, StreamMessage, SubscribeRequest } from './gen/chat_pb';

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

  // This useEffect hook subscribes to the message stream when the user is joined
  useEffect(() => {
    if (isJoined && username) {
      const controller = new AbortController();

      async function subscribeToMessages() {
        // The 'SubscribeRequest' class is now correctly imported
        const request = new SubscribeRequest({ user: username });
        try {
          const stream = client.subscribe(request, { signal: controller.signal });
          for await (const res of stream) {
            if (res.event.case === "message") {
              setMessages(prev => [...prev, res.event.value]);
            }
          }
        } catch (err) {
          if (err.name !== 'AbortError') {
            console.error("Subscription failed:", err);
            alert("Connection to the server was lost.");
            setIsJoined(false);
            setMessages([]);
          }
        }
      }

      subscribeToMessages();

      return () => {
        controller.abort();
      };
    }
  }, [isJoined, username]);

  const handleSendMessage = async (messageText: string) => {
    if (!username) return;

    // The 'ChatMessage' class is now correctly imported
    const request = new ChatMessage({
      user: username,
      message: messageText,
    });

    try {
      await client.sendMessage(request);
    } catch (err) {
      console.error("Failed to send message:", err);
      alert("Could not send message.");
    }
  };

  return (
    <main className="bg-gray-100">
      {isJoined ? (
        <ChatScreen username={username} messages={messages.map(m => m.toJson()) as any} onSendMessage={handleSendMessage} />
      ) : (
        <JoinScreen onJoin={handleJoin} />
      )}
    </main>
  );
}

export default App;