import React, { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { AiOutlineCopy, AiOutlineLink, AiOutlineFile, AiOutlineFileImage, AiOutlinePlus, AiOutlineDelete } from "react-icons/ai";
import { MdContentPaste } from "react-icons/md";
import { GetAllClipBoardItems, AddClipBoardItem, DeleteClipBoardItem, GetClipboardContent } from "../wailsjs/go/main/App";
import { ClipboardItemDbRow } from './types/clipboard';
import axios from 'axios';
import * as cheerio from 'cheerio';

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
        item.content.toLowerCase().trim().includes(searchQuery.trim().toLowerCase())
    );

    const fetchLinkMetadata = async (url: string): Promise<LinkMetadata> => {
        try {
            const { data: html } = await axios.get(url);

            const $ = cheerio.load(html);

            const title = $('head > title').text() || url;
            const description = $('meta[name="description"]').attr('content') || ''
            let favicon = $('link[rel="icon"]').attr('href') || $('link[rel="shortcut icon"]').attr('href');

            if (favicon && !favicon.startsWith('http')) {
                const urlObject = new URL(url);
                favicon = `${urlObject.origin}${favicon}`;
            } else if (!favicon) {
                // Provide a default favicon or handle no favicon scenario
                favicon = "No favicon found";
            }

            return {
                title,
                description,
                favicon
            };
        } catch (error) {
            console.error(`Error fetching metadata for URL ${url}:`, error);
            throw new Error(`Unable to fetch metadata for the provided URL`);
        }
    };

    const renderClipboardItem = (item: ClipboardItemDbRow) => {
        if (item.type.startsWith('image/') || (item.type === 'link' && item.content.match(/\.(jpeg|jpg|gif|png|media)$/))) {
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
            return <p className="break-all text-sm line-clamp-2 bg-gray-700 p-1 rounded-md">{item.content}</p>;
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
        <div className="bg-[#282828] text-[#ebdbb2] min-h-screen p-4">
            <div className="max-w-2xl mx-auto">
                <h1 className="text-2xl font-bold mb-4 text-[#fabd2f]">Voxel</h1>

                <div className="mb-4 flex space-x-2">
                    <Input
                        placeholder="Search clipboard items"
                        value={searchQuery}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
                        className="bg-[#3c3836] border-[#504945] text-[#ebdbb2] flex-grow"
                    />
                </div>

                <div className="mb-4 flex space-x-2">
                    <Input
                        placeholder="New clipboard item"
                        value={newItemContent}
                        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewItemContent(e.target.value)}
                        className="bg-[#3c3836] border-[#504945] text-[#ebdbb2] flex-grow"
                    />
                    <Button onClick={pasteFromClipboard} className="bg-[#98971a] hover:bg-[#b8bb26]">
                        <MdContentPaste className="mr-2" /> Paste
                    </Button>
                    <Button onClick={addNewItem} className="bg-[#458588] hover:bg-[#83a598]">
                        <AiOutlinePlus className="mr-2" /> Add
                    </Button>
                </div>

                <ScrollArea className="h-[calc(100vh-220px)]">
                    <div className="space-y-2">
                        {filteredItems.map((item) => (
                            <div key={item.id} className="bg-[#3c3836] rounded-md p-3 flex items-center space-x-3">
                                <div>
                                    {item.type.startsWith('image/') ? (
                                        <AiOutlineFileImage className="text-xl text-[#458588] flex-shrink-0" />
                                    ) : item.type === 'link' ? (
                                        <AiOutlineLink className="text-xl text-[#98971a] flex-shrink-0" />
                                    ) : (
                                        <AiOutlineFile className="text-xl text-[#7c6f64] flex-shrink-0" />
                                    )}
                                </div>
                                <div className="flex flex-wrap min-w-0 flex-1">
                                    {renderClipboardItem(item)}
                                </div>
                                <div>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        onClick={() => copyToClipboard(item.content)}
                                        className="text-[#ebdbb2] hover:text-[#fe8019] flex-shrink-0"
                                    >
                                        <AiOutlineCopy
                                            className="h-6 w-6" />
                                    </Button>
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        onClick={() => deleteItem(item.id)}
                                        className="text-[#cc241d] hover:text-[#fb4934] flex-shrink-0"
                                    >
                                        <AiOutlineDelete className="h-6 w-6" />
                                    </Button>
                                </div>
                            </div>
                        ))}
                    </div>
                </ScrollArea>
            </div>
        </div>
    );
};

export default ClipboardManager;
