import React from "react";
import { ffBlurProMedium } from "~/assets/fonts";

type Props = React.ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: 'primary' | 'secondary'
}

export default function Button({ children, className, variant = 'primary', ...props }: Props) {
    const variantStyles = variant === 'primary' 
        ? `${props.disabled ? 'bg-theme-red-shadow text-gray-300' : 'bg-theme-red text-white'} active:bg-theme-red-shadow`
        : 'bg-theme-white text-theme-red border border-theme-red'
    
    return (
        <button className={`${className} ${ffBlurProMedium.className} ${variantStyles} p-2 rounded-md transition-colors duration-100`} {...props}>
            {children}
        </button>
    )
}