

export interface ClipboardItem {
    content: string;
    timestamp ?: string;
    type: string;
    categories: string;
}

export interface ClipboardItemDbRow extends ClipboardItem {
    id: number;
}

