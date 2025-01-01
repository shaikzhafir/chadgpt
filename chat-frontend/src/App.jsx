import { useEffect, useState } from 'react';
import axios from 'axios';
import { ModelSelector } from './components/ModelSelector';

function App() {
  const [model, setModel] = useState('openai');
  const [sessions, setSessions] = useState([]);
  const [currentSessionId, setCurrentSessionId] = useState(null);
  const [currentMessages, setCurrentMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');

  // create a new session on init and cleanup 
  useEffect(() => {
    // Check if there's no current session
    if (!currentSessionId || !sessions.length) {
      createNewSession();
    }
    return () => {
      console.log('Cleaning up');
    };
  }, []); // Empty dependency array for init only


  const createNewSession = () => {
    console.log('Creating new session');
    const newSessionId = `session-${Date.now()}`; // You could generate this based on some logic or from the backend
    setSessions([...sessions, { sessionId: newSessionId, messages: [] }]);
    setCurrentSessionId(newSessionId);
    setCurrentMessages([]);
  };

  const switchSession = (sessionId) => {
    const session = sessions.find((s) => s.sessionId === sessionId);
    if (session) {
      setCurrentSessionId(sessionId);
      setCurrentMessages(session.messages);
    }
  };

  const sendMessage = async () => {
    if (inputMessage.trim() === '') return;

    try {
      console.log('Sending message with session-id', currentSessionId);
      const url = `/chat-${model}`; 
      console.log('url:', url);
      const response = await axios.post(url, { message: inputMessage }, {
        headers: {
          'session-id': currentSessionId,
        },
      });

      const newMessage = { user: 'user', text: inputMessage };
      const newMessages = [...currentMessages, newMessage, { user: 'bot', text: response.data }];

      setCurrentMessages(newMessages);

      // Update the session in the state
      setSessions(sessions.map(s => 
        s.sessionId === currentSessionId ? { ...s, messages: newMessages } : s
      ));

      setInputMessage(''); // clear the input
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  return (
    <div className="flex h-screen">
      {/* Sidebar: List of Sessions */}
      <div className="w-1/4 bg-gray-100 p-4">
        <button
          className="bg-blue-500 text-white py-2 px-4 rounded mb-4"
          onClick={createNewSession}
        >
          New Session
        </button>
        <ModelSelector
        selectedModel={model}
        onModelChange={setModel}
        />
        <ul>
          {sessions.map((session) => (
            <li key={session.sessionId}>
              <button
                className={`block w-full text-left p-2 ${
                  currentSessionId === session.sessionId ? 'bg-blue-200' : 'bg-gray-200'
                }`}
                onClick={() => switchSession(session.sessionId)}
              >
                {session.sessionId}
              </button>
            </li>
          ))}
        </ul>
      </div>

      {/* Chat Window */}
      <div className="w-3/4 flex flex-col bg-white p-4">
        <div className="flex-1 overflow-auto">
          {currentMessages.map((msg, index) => (
            <div key={index} className={`mb-2 ${msg.user === 'user' ? 'text-right' : 'text-left'}`}>
              <span
                className={`inline-block p-2 rounded ${
                  msg.user === 'user' ? 'bg-blue-300' : 'bg-gray-300'
                }`}
              >
                {msg.text}
              </span>
            </div>
          ))}
        </div>
        <div className="mt-4 flex">
          <input
            type="text"
            className="flex-1 border p-2 rounded mr-2"
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            placeholder="Type your message..."
          />
          <button
            className="bg-blue-500 text-white py-2 px-4 rounded"
            onClick={sendMessage}
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}

export default App
