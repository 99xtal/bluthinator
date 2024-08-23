import React from "react";
import { ffBlurProMedium } from "~/fonts";

type Props = React.ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: 'primary' | 'secondary'
}

export default function Button({ children, className, variant = 'primary', ...props }: Props) {
    const variantStyles = variant === 'primary' ? 'bg-theme-red text-white' : 'bg-theme-white text-theme-red border border-theme-red'
    return (
        <button className={`${className} ${ffBlurProMedium.className} ${variantStyles} p-2 rounded-md`} {...props}>
            {children}
        </button>
    )
}