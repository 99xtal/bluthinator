import Link from "next/link";
import React from "react";
import { ffBlurProMedium } from "~/assets/fonts";

export default function TextLink({ children, href }: { children: string, href: string }) {
    return (
        <Link href={href} className={`${ffBlurProMedium.className} text-md text-theme-black hover:underline`}>
            <p>{children}</p>
        </Link>
    )
}