import React from "react";
import { ffBlurProMedium } from "~/assets/fonts";

interface Props {
    fontSize?: number
    children: string
}

export default function Caption({ children, fontSize = 36 }: Props) {
    return (
        <h2 className={`${ffBlurProMedium.className} text-theme-yellow`} style={{ textShadow: '2px 2px 4px rgba(0, 0, 0, 0.9)', fontSize: `${fontSize}px`, lineHeight: `${fontSize}px` }}>
            {children}
        </h2>
    )
}