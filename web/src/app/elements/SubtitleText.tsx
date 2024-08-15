import React from "react";
import { defonteRegular } from "~/fonts";

export default function SubtitleText({ children }: { children: string }) {
    return (
        <h3 className={`${defonteRegular.className} text-xl text-theme-red`}>{children}</h3>
    )
}