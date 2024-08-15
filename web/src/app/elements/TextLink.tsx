import Link from "next/link";
import React from "react";
import { defonteRegular } from "~/fonts";

export default function TextLink({ children, href }: { children: string, href: string }) {
    return (
        <Link href={href} className={`${defonteRegular.className} text-md text-theme-black hover:underline`}>
            <p>{children}</p>
        </Link>
    )
}