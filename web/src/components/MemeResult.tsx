'use client';

import Link from "next/link";
import Image from "next/image";
import { useSearchParams } from "next/navigation";

interface Props {
    episode: string;
    timestamp: string;
    caption?: string;
}

export default function MemeResult({ episode, timestamp, caption: captionProp }: Props) {
    const searchParams = useSearchParams();
    const caption = captionProp ?? searchParams.get('b')?.toString();
    
    let url = `${process.env.NEXT_PUBLIC_API_HOST}/caption/${episode}/${timestamp}`;
    if (caption) {
        url += `?b=${caption}`;
    }

    const imageUrl = caption 
        ? `/img/caption/${episode}/${timestamp}/${caption}` 
        : `/img/${episode}/${timestamp}/large.jpg`;
 

    return (
        <Link href={imageUrl}>
            <Image src={url} alt="caption result" width={640} height={360} className="outline outline-4 outline-theme-black" />
        </Link>
    )
}