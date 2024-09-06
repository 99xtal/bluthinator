'use client';

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useDebouncedCallback } from 'use-debounce';
import { logEvent } from "~/utils/firebase";

export default function Search({ placeholder, className }: { placeholder: string, className?: string }) {
    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { push, replace } = useRouter();  

    const handleSearch = useDebouncedCallback((query: string) => {
        logEvent("search", { query });
        const params = new URLSearchParams(searchParams)
        if (query) {
            params.set("q", query);
        } else {
            params.delete("q");
        }
        const path = `/?${params.toString()}`
        if (pathname !== "/") {
            push(path);
        } else {
            replace(path);
        }
    }, 300);

    return (
        <input 
            placeholder={placeholder} 
            onChange={(e) => {
                handleSearch(e.target.value);
            }} 
            defaultValue={searchParams.get('q')?.toString()}
            className={className}
        />
    )
}