import React from "react";
import { ffBlurProMedium } from "~/fonts";

export default function Caption({ children }: { children: string }) {
    return (
        <h2 className={`${ffBlurProMedium.className} text-4xl text-theme-yellow`} style={{ textShadow: '2px 2px 4px rgba(0, 0, 0, 0.9)' }}>
            {children}
        </h2>
    )
}