
import { useState, useEffect } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { Greet, GetAllClipBoardItems, AddClipBoardItem } from "../wailsjs/go/main/App";
import { ClipboardItemDbRow } from './types/clipboard'; // Import the TypeScript types
import { AiFillFire } from "react-icons/ai"

function App() {
    const [resultText, setResultText] = useState("Please enter content for the clipboard below 👇");
    const [content, setContent] = useState(''); // Renamed 'name' to 'content'
    const [clipboardItems, setClipboardItems] = useState<ClipboardItemDbRow[]>([]); // State to hold clipboard items

    const updateResultText = (result: string) => setResultText(result);
    const updateContent = (e: any) => setContent(e.target.value);

    async function greet() {
        const string = await Greet(content)
        updateResultText(string);
        await AddClipBoardItem(content, ["default", "default"], "text/plain")
        updateResultText("Added to clipboard");
        setContent('');
        await fetchClipboardItems()
    }

    // Function to fetch all clipboard items
    const fetchClipboardItems = async () => {
        try {
            const items = await GetAllClipBoardItems();
            setClipboardItems(items);
            console.log("Clipboard items fetched:", items);
        } catch (err) {
            console.error("Error fetching clipboard items:", err);
        }
    };

    useEffect(() => {
        fetchClipboardItems();
    }, []);

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo" />
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input
                    id="content"
                    className="input"
                    value={content} // Controlled input
                    onChange={updateContent}
                    autoComplete="off"
                    name="input"
                    type="text"
                    placeholder="Enter clipboard content"
                />
                <button className="btn" onClick={greet}>Add </button>
            </div>

            {/* Render the clipboard items */}
            <div id="clipboard-list">
                <h3>Clipboard Items</h3>
                <ul>
                    {clipboardItems.map(item => (
                        <li key={item.id}>
                            <AiFillFire className="icon" />
                            <p><strong>Content:</strong> {item.content}</p>
                            <p><strong>Type:</strong> {item.type}</p>
                            <p><strong>Categories:</strong> {item.categories}</p>
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    );
}

export default App;

