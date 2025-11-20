// frontend/src/components/ChatScreen.tsx
import { useState, useRef, useEffect } from "react";
// Import the TYPE ONLY. This avoids bundling the whole module.
import type { ChatMessage as ChatMessageType } from "../grpc/chat_pb.d";

interface ChatScreenProps {
  username: string;
  messages: ChatMessageType.AsObject[];
  onSendMessage: (message: string) => void;
}

export function ChatScreen({ username, messages, onSendMessage }: ChatScreenProps) {
  const [message, setMessage] = useState("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const handleSend = () => {
    if (message.trim()) {
      onSendMessage(message.trim());
      setMessage("");
    }
  };

  // Auto-scroll to the latest message
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <div className="flex flex-col h-screen p-2 sm:p-4 bg-gray-100">
      <div className="flex flex-col flex-grow bg-white rounded-lg shadow-xl overflow-hidden">
        {/* Header */}
        <div className="p-4 border-b bg-blue-600 text-white font-semibold text-lg text-center">
          You are logged in as: {username}
        </div>

        {/* Messages Area */}
        <div className="flex-grow p-4 overflow-y-auto">
          <div className="space-y-4">
            {messages.map((msg, index) => {
              const isCurrentUser = msg.user === username;
              return (
                <div
                  key={index}
                  className={`flex items-end gap-3 ${isCurrentUser ? "justify-end" : "justify-start"}`}
                >
                  {/* Avatar for others */}
                  {!isCurrentUser && (
                    <div className="w-8 h-8 rounded-full bg-gray-300 flex items-center justify-center text-gray-600 font-bold flex-shrink-0">
                      {msg.user.substring(0, 2).toUpperCase()}
                    </div>
                  )}

                  {/* Message Bubble */}
                  <div
                    className={`rounded-2xl p-3 max-w-xs md:max-w-md ${
                      isCurrentUser
                        ? "bg-blue-600 text-white rounded-br-none"
                        : "bg-gray-200 text-gray-800 rounded-bl-none"
                    }`}
                  >
                    {!isCurrentUser && <p className="text-sm font-bold mb-1 text-gray-600">{msg.user}</p>}
                    <p className="text-base">{msg.message}</p>
                  </div>
                </div>
              );
            })}
            {/* Anchor for scrolling */}
            <div ref={messagesEndRef} />
          </div>
        </div>

        {/* Input Form */}
        <div className="p-4 border-t bg-gray-50">
          <div className="flex items-center gap-2">
            <input
              type="text"
              placeholder="Type your message..."
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              onKeyPress={(e) => e.key === "Enter" && handleSend()}
              className="flex-grow px-4 py-2 text-base text-gray-700 bg-white border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
            />
            <button
              onClick={handleSend}
              className="bg-blue-600 text-white px-5 py-2 rounded-full hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition flex-shrink-0"
            >
              Send
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}