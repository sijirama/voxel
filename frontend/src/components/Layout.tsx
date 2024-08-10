
import { ReactNode } from 'react'
import Sidebar from './custom/sidebar'

interface Props {
    children: ReactNode
}

export default function Layout({ children }: Props) {
    return (
        <main className='flex gap-2'>
            <Sidebar />
            {children}
        </main>
    )
}

