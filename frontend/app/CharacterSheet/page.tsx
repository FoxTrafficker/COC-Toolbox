'use client'
import {useState, useEffect} from 'react';

export default function Home() {
    const [jsonMessage, setJsonMessage] = useState('');
    const [newMessage, setNewMessage] = useState('');
    const [socket, setSocket] = useState<WebSocket | null>(null);

    // 初始化 WebSocket 连接并获取 JSON 文件内容
    useEffect(() => {
        const ws = new WebSocket('ws://localhost:8080/ws');

        ws.onopen = () => {
            console.log("WebSocket 连接已建立");
        };

        // 接收来自服务器的消息
        ws.onmessage = (event) => {
            console.log("收到服务器推送的消息:", event.data);
            setJsonMessage(event.data);
        };

        ws.onclose = () => {
            console.log("WebSocket 连接已关闭");
        };

        ws.onerror = (error) => {
            console.error("WebSocket 错误:", error);
        };

        setSocket(ws);

        return () => {
            if (ws) ws.close();  // 组件卸载时关闭 WebSocket 连接
        };
    }, []);

    // 通过 WebSocket 发送消息更新
    const updateJsonMessage = () => {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(newMessage);  // 通过 WebSocket 发送消息
            setNewMessage('');        // 清空输入框
        }
    };

    return (
        <div>
            <h1>WebSocket 客户端和 JSON 处理</h1>
            <h2>当前 JSON 文件的消息: {jsonMessage}</h2>

            <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder="输入新的消息"
            />
            <button onClick={updateJsonMessage}>更新 JSON 消息</button>
        </div>
    );
}
