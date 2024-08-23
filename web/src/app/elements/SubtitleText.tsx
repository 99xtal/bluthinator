import React from "react";
import { ffBlurProMedium } from "~/fonts";

export default function SubtitleText({ children }: { children: string }) {
    return (
        <h3 className={`${ffBlurProMedium.className} text-xl text-theme-red`}>{children}</h3>
    )
}