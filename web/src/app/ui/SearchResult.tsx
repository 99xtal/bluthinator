import Image from 'next/image';
import Link from "next/link";

export default function SearchResult({ result }: { result: any }) {
    const getImgUrl = (result: any) => {
        return `${process.env.NEXT_PUBLIC_IMG_HOST}/frames/${result.episode}/${result.timestamp}/small.png`;
    }

    return (
        <Link href={`/episode/${result.episode}`}>
            <Image 
                src={getImgUrl(result)} 
                alt={`${result.episode}: ${result.timestamp}`} 
                width={400} 
                height={240} 
                className="box-border hover:outline hover:outline-8 hover:outline-theme-black" 
            />
        </Link>
    )
}