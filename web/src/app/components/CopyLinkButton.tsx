"use client";

import Image from "next/image";
import { Button } from "../elements";
import copySvg from "~/assets/svg/copy.svg";
import { useState } from "react";
import { ffBlurProMedium } from "~/fonts";

export default function CopyLinkButton() {
    const [copied, setCopied] = useState(false);

    const handleClick = async () => {
        await navigator.clipboard.writeText(window.location.href);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    }

    return (
        <div className="relative flex flex-col items-center">
            <Button onClick={handleClick} title="Copy Link" className="px-6 py-2">
                <Image src={copySvg} alt="Copy Link" width={32} height={32} />
            </Button>
            {copied && <span className={`absolute top-full text-sm text-theme-black mt-2 ${ffBlurProMedium.className}`}>Link Copied!</span>}
        </div>
    )
}