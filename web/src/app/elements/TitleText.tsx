import React from "react";
import { ffBlurProMedium } from "~/fonts";

export default function TitleText({ children }: { children: string }) {
    return (
        <h1 className={`${ffBlurProMedium.className} text-3xl`}>{children}</h1>
    )
}