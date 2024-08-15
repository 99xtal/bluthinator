import React from "react";
import { defonteRegular } from "~/fonts";

export default function TitleText({ children }: { children: string }) {
    return (
        <h1 className={`${defonteRegular.className} text-3xl`}>{children}</h1>
    )
}