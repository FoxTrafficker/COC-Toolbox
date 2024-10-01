'use client'

import "./globals.css";
import {Inter} from "next/font/google";
import React from "react";
import SideBar from "@/components/SideBar";

const inter = Inter({subsets: ["latin"]});


export default function RootLayout({children}: { children: React.ReactNode }) {
    const [drawerOpen, setDrawerOpen] = React.useState(true);

    return (
        <html lang="en">
            <body className={inter.className}>
                <header/>

                <SideBar drawerOpen={drawerOpen} setDrawerOpen={setDrawerOpen}/>

                <main style={{marginLeft: drawerOpen ? '240px' : '0px', transition: 'margin-left 0.35s'}}>
                    {children}
                </main>
            </body>
        </html>
    );
}