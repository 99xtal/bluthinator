import React from "react";
import { defonteRegular } from "~/fonts";

type Props = React.ButtonHTMLAttributes<HTMLButtonElement>

export default function Button({ children, ...props }: Props) {
    return (
        <button className={`${defonteRegular.className} flex-1 bg-theme-red text-white p-2 rounded-md`} {...props}>
            {children}
        </button>
    )
}