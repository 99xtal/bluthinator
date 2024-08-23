import React from "react";
import { ffBlurProMedium } from "~/fonts";

type Props = React.ButtonHTMLAttributes<HTMLButtonElement>

export default function Button({ children, className, ...props }: Props) {
    return (
        <button className={`${className} ${ffBlurProMedium.className} flex-1 bg-theme-red text-white p-2 rounded-md`} {...props}>
            {children}
        </button>
    )
}