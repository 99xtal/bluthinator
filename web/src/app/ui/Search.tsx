'use client';

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useDebouncedCallback } from 'use-debounce';

export default function Search({ placeholder }: { placeholder: string}) {
    const searchParams = useSearchParams();
    const pathname = usePathname();
    const { replace } = useRouter();  

    const handleSearch = useDebouncedCallback((query: string) => {
        const params = new URLSearchParams(searchParams)
        if (query) {
            params.set("q", query);
        } else {
            params.delete("q");
        }
        replace(`${pathname}?${params.toString()}`);
    }, 500);

    return (
        <input 
            placeholder={placeholder} 
            onChange={(e) => {
                handleSearch(e.target.value);
            }} 
            defaultValue={searchParams.get('q')?.toString()}
        />
    )
}