// frontend/src/components/BoardView.tsx
import React, { useState, useEffect } from 'react';
import useWebSocket from 'react-use-websocket';

interface MessagePayload {
    username: string;
    message: string;
}

export const BoardView = () => {
    const [messageHistory, setMessageHistory] = useState<MessagePayload[]>([]);
    const [currentMessage, setCurrentMessage] = useState('');

    const token = localStorage.getItem('authToken');
    const socketUrl = `ws://localhost:8082/ws/board/general?token=${token}`;

    const { sendMessage, lastMessage } = useWebSocket(socketUrl, {
        onOpen: () => console.log('WebSocket connection opened.'),
        onClose: () => console.log('WebSocket connection closed.'),
        shouldReconnect: (closeEvent) => true,
    });

    useEffect(() => {
        if (lastMessage !== null) {
            const payload = JSON.parse(lastMessage.data);
            setMessageHistory((prev) => [...prev, payload]);
        }
    }, [lastMessage]);

    const handleSendMessage = () => {
        if (currentMessage.trim()) {
            sendMessage(currentMessage);
            setCurrentMessage('');
        }
    };

    return (
        <div className="board-container">
            <div className="message-history">
                {messageHistory.map((msg, idx) => (
                    <div key={idx} className="message">
                        <strong>{msg.username}:</strong> {msg.message}
                    </div>
                ))}
            </div>
            <div className="message-input">
                <input
                    type="text"
                    value={currentMessage}
                    onChange={(e) => setCurrentMessage(e.target.value)}
                    onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                    placeholder="Type your message..."
                />
                <button onClick={handleSendMessage}>Send</button>
            </div>
        </div>
    );
};