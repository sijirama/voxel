import React, { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AiOutlineCopy, AiOutlineLink, AiOutlineFile, AiOutlineFileImage, AiOutlinePlus, AiOutlineDelete } from "react-icons/ai";
import { MdContentPaste } from "react-icons/md";
import { GetAllClipBoardItems, AddClipBoardItem, DeleteClipBoardItem, GetClipboardContent } from "../wailsjs/go/main/App";
import { ClipboardItemDbRow } from './types/clipboard';

interface LinkMetadata {
    title: string | null;
    description: string | null;
    favicon: string | null;
}

const ClipboardManager: React.FC = () => {
    const [clipboardItems, setClipboardItems] = useState<ClipboardItemDbRow[]>([]);
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [newItemContent, setNewItemContent] = useState<string>('');

    useEffect(() => {
        fetchClipboardItems();
    }, []);

    const fetchClipboardItems = async () => {
        try {
            const items = await GetAllClipBoardItems();
            setClipboardItems(items);
        } catch (err) {
            console.error("Error fetching clipboard items:", err);
        }
    };

    const addNewItem = async () => {
        if (newItemContent.trim()) {
            const type = newItemContent.startsWith('http') ? 'link' : 'text/plain';
            await AddClipBoardItem(newItemContent, ["default"], type);
            setNewItemContent('');
            fetchClipboardItems();
        }
    };

    const pasteFromClipboard = async () => {
        try {
            const content = await GetClipboardContent();
            setNewItemContent(content);
        } catch (err) {
            console.error("Error pasting from clipboard:", err);
        }
    };

    const copyToClipboard = (content: string) => {
        navigator.clipboard.writeText(content);
    };

    const deleteItem = async (id: number) => {
        try {
            await DeleteClipBoardItem(id);
            fetchClipboardItems();
        } catch (err) {
            console.error("Error deleting clipboard item:", err);
        }
    };

    const filteredItems = clipboardItems.filter(item =>
        item.content.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const fetchLinkMetadata = async (url: string): Promise<LinkMetadata> => {
        // In a real application, you'd want to implement this on the backend
        // For demonstration, we'll just return a mock object
        return {
            title: "Sample Title",
            description: "This is a sample description for the link.",
            favicon: "https://example.com/favicon.ico"
        };
    };

    const renderClipboardItem = (item: ClipboardItemDbRow) => {
        if (item.type.startsWith('image/') || (item.type === 'link' && item.content.match(/\.(jpeg|jpg|gif|png)$/))) {
            return (
                <div className="relative w-full h-40">
                    <img
                        src={item.content}
                        alt="Clipboard content"
                        className="w-full h-full object-cover rounded-lg"
                        onError={(e: React.SyntheticEvent<HTMLImageElement, Event>) => {
                            e.currentTarget.src = "/api/placeholder/400/320";
                            e.currentTarget.alt = "Failed to load image";
                        }}
                    />
                </div>
            );
        } else if (item.type === 'link') {
            return <LinkItem url={item.content} />;
        } else {
            return <p className="break-all line-clamp-2">{item.content}</p>;
        }
    };

    const LinkItem: React.FC<{ url: string }> = ({ url }) => {
        const [metadata, setMetadata] = useState<LinkMetadata | null>(null);

        useEffect(() => {
            fetchLinkMetadata(url).then(setMetadata);
        }, [url]);

        return (
            <div className="flex items-center space-x-3">
                {metadata?.favicon && (
                    <img src={metadata.favicon} alt="Site favicon" className="w-6 h-6" />
                )}
                <div>
                    <a href={url} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline font-medium">
                        {metadata?.title || url}
                    </a>
                    {metadata?.description && (
                        <p className="text-sm text-gray-400 line-clamp-2">{metadata.description}</p>
                    )}
                </div>
            </div>
        );
    };

    return (
        <div className="bg-gray-900 text-white min-h-screen p-4">
            <div className="max-w-2xl mx-auto">
                <h1 className="text-2xl font-bold mb-4">Clipboard History</h1>

                <div className="mb-4 flex space-x-2">
                    <Input
                        placeholder="Search clipboard items"
                        value={searchQuery}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
                        className="bg-gray-800 border-gray-700 text-white flex-grow"
                    />
                </div>

                <div className="mb-4 flex space-x-2">
                    <Input
                        placeholder="New clipboard item"
                        value={newItemContent}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewItemContent(e.target.value)}
                        className="bg-gray-800 border-gray-700 text-white flex-grow"
                    />
                    <Button onClick={pasteFromClipboard} className="bg-green-600 hover:bg-green-700">
                        <MdContentPaste className="mr-2" /> Paste
                    </Button>
                    <Button onClick={addNewItem} className="bg-blue-600 hover:bg-blue-700">
                        <AiOutlinePlus className="mr-2" /> Add
                    </Button>
                </div>

                <ScrollArea className="h-[calc(100vh-220px)]">
                    <div className="space-y-2">
                        {filteredItems.map((item) => (
                            <div key={item.id} className="bg-gray-800 rounded-lg p-3 flex items-center space-x-3">
                                {item.type.startsWith('image/') ? (
                                    <AiOutlineFileImage className="text-xl text-blue-500 flex-shrink-0" />
                                ) : item.type === 'link' ? (
                                    <AiOutlineLink className="text-xl text-green-500 flex-shrink-0" />
                                ) : (
                                    <AiOutlineFile className="text-xl text-gray-500 flex-shrink-0" />
                                )}
                                <div className="flex-grow min-w-0">
                                    {renderClipboardItem(item)}
                                </div>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => copyToClipboard(item.content)}
                                    className="text-gray-400 hover:text-white flex-shrink-0"
                                >
                                    <AiOutlineCopy className="h-4 w-4" />
                                </Button>
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => deleteItem(item.id)}
                                    className="text-red-400 hover:text-red-600 flex-shrink-0"
                                >
                                    <AiOutlineDelete className="h-4 w-4" />
                                </Button>
                            </div>
                        ))}
                    </div>
                </ScrollArea>
            </div>
        </div>
    );
};

export default ClipboardManager;
