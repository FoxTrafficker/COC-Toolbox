"use client"
import {useState, useEffect, useRef} from "react";

interface Character {
    name: string;
    agility: number;
    avatar: string;
}

interface WebSocketMessage {
    type: string;
    payload: string;
}

const CharacterList = ({characters, deleteCharacter}: { characters: Character[]; deleteCharacter: (index: number) => void }) => {
    return (
        <div style={styles.indicator}>
            {characters.map((character, index) => (
                <div key={index} style={styles.character}>
                    <img
                        src={"http://localhost:8000" + character.avatar}
                        alt={character.name}
                        style={styles.avatar}
                    />
                    <p>{character.name}</p>
                    <button onClick={() => deleteCharacter(index)}>删除</button>
                </div>
            ))}
        </div>
    );
};

const NextTurnButton = ({nextTurn}: { nextTurn: () => void }) => {
    return (
        <button style={styles.nextTurnButton} onClick={nextTurn}>
            下一回合
        </button>
    );
};

const AddCharacterForm = ({newName, setNewName, newAgility, setNewAgility, addCharacter}: {
    newName: string;
    setNewName: (name: string) => void;
    newAgility: string;
    setNewAgility: (agility: string) => void;
    addCharacter: () => void;
}) => {
    return (
        <div style={styles.addForm}>
            <input
                type="text"
                placeholder="角色名称"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
            />
            <input
                type="number"
                placeholder="敏捷度"
                value={newAgility}
                onChange={(e) => setNewAgility(e.target.value)}
            />
            <button onClick={addCharacter}>添加角色</button>
        </div>
    );
};

const Home: React.FC = () => {
    const [characters, setCharacters] = useState<Character[]>([]);
    const [newName, setNewName] = useState<string>("");
    const [newAgility, setNewAgility] = useState<string>("");
    const ws = useRef<WebSocket | null>(null);


    useEffect(() => {
        // 连接 WebSocket 服务器
        ws.current = new WebSocket("ws://localhost:8000/ws");

        // 接收消息
        ws.current.onmessage = (event) => {
            const message: WebSocketMessage = JSON.parse(event.data);
            if (message.type === "INITIAL_STATE" || message.type === "UPDATE_CHARACTERS") {
                setCharacters(JSON.parse(message.payload));
            }
        };


        return () => {if (ws.current) {ws.current.close()}};

    }, []);

    // 发送 WebSocket 消息
    const sendMessage = (message: WebSocketMessage) => {
        if (ws.current && ws.current.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(message));
        }
    };

    // 下一回合函数
    const nextTurn = () => {
        const updatedCharacters = [...characters];
        const firstCharacter = updatedCharacters.shift();
        if (firstCharacter) updatedCharacters.push(firstCharacter);
        setCharacters(updatedCharacters);
        sendMessage({
            type: "UPDATE_CHARACTERS",
            payload: JSON.stringify(updatedCharacters)
        });
    };

    // 添加角色函数
    const addCharacter = () => {
        if (!newName || !newAgility) return;
        const newCharacter: Character = {
            name: newName,
            agility: parseInt(newAgility),
            avatar: `/avatars/default.webp`,
        };
        const updatedCharacters = [...characters, newCharacter].sort((a, b) => b.agility - a.agility);
        setCharacters(updatedCharacters);
        setNewName("");
        setNewAgility("");
        sendMessage({
            type: "UPDATE_CHARACTERS",
            payload: JSON.stringify(updatedCharacters)
        });
    };

    // 删除角色函数
    const deleteCharacter = (index: number) => {
        const updatedCharacters = characters.filter((_, i) => i !== index);
        setCharacters(updatedCharacters);
        sendMessage({
            type: "UPDATE_CHARACTERS",
            payload: JSON.stringify(updatedCharacters)
        });
    };

    const resetCharacters = () => {
        const sortedCharacters = [...characters].sort((a, b) => b.agility - a.agility);
        setCharacters(sortedCharacters);
        sendMessage({
            type: "RESET_CHARACTERS",
            payload: JSON.stringify(sortedCharacters),
        });
    };

    return (
        <div style={styles.container}>
            <h1>COC 战斗轮指示器</h1>
            <CharacterList characters={characters} deleteCharacter={deleteCharacter}/>
            <NextTurnButton nextTurn={nextTurn}/>
            <button style={styles.resetButton} onClick={resetCharacters}>
                重置战斗指示器
            </button>
            <h2/>
            <AddCharacterForm
                newName={newName}
                setNewName={setNewName}
                newAgility={newAgility}
                setNewAgility={setNewAgility}
                addCharacter={addCharacter}
            />
        </div>
    );
};

export default Home;

const styles = {
    container: {
        textAlign: "center" as const,
        padding: "20px",
    },
    indicator: {
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        margin: "20px auto",
        width: "80%",
        border: "2px solid #000",
        padding: "10px",
        backgroundColor: "#f4f4f4",
    },
    character: {
        display: "flex",
        flexDirection: "column" as const,
        alignItems: "center",
        margin: "0 10px",
    },
    avatar: {
        width: "50px",
        height: "50px",
        borderRadius: "50%",
        objectFit: "cover" as const,
        border: "2px solid #000",
    },
    nextTurnButton: {
        padding: "10px 20px",
        fontSize: "16px",
        backgroundColor: "#4caf50",
        color: "white",
        border: "none",
        borderRadius: "5px",
        cursor: "pointer",
    },
    addForm: {
        margin: "20px auto",
    },
    resetButton: {
        padding: "10px 20px",
        fontSize: "16px",
        backgroundColor: "#f44336",
        color: "white",
        border: "none",
        borderRadius: "5px",
        cursor: "pointer",
        marginTop: "10px",
    },
};
