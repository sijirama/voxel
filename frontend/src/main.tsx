import React from 'react'
import { createRoot } from 'react-dom/client'
import './style.css'
import App from './App'
import Layout from './components/Layout'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <App />
        <Layout>
        </Layout>
    </React.StrictMode>
)
