// frontend/src/components/JoinScreen.tsx
import { useState } from "react";

interface JoinScreenProps {
  onJoin: (username: string) => void;
}

export function JoinScreen({ onJoin }: JoinScreenProps) {
  const [username, setUsername] = useState("");

  const handleJoin = () => {
    if (username.trim()) {
      onJoin(username.trim());
    }
  };

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <div className="w-full max-w-sm p-8 space-y-6 bg-white shadow-lg rounded-xl">
        <h2 className="text-3xl font-bold text-center text-gray-800">
          Welcome to gRPC Chat
        </h2>
        <div className="space-y-4">
          <div>
            <label htmlFor="username" className="sr-only">
              Username
            </label>
            <input
              type="text"
              id="username"
              placeholder="Enter your username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              onKeyPress={(e) => e.key === "Enter" && handleJoin()}
              className="w-full px-4 py-3 text-lg text-gray-700 bg-gray-50 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
            />
          </div>
          <button
            onClick={handleJoin}
            className="w-full px-4 py-3 text-lg font-semibold text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition"
          >
            Join Chat
          </button>
        </div>
      </div>
    </div>
  );
}