"use client";

import { PropsWithChildren } from "react";
import { useRouter } from "next/navigation";

type Props = PropsWithChildren & {
    className?: string;
}

export default function GoBackLink({ className, children }: Props) {
    const router = useRouter();

    return (
        <button onClick={() => router.back()} className={className}>
            {children}
        </button>
    )
}