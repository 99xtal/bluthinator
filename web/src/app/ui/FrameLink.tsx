import Image from 'next/image';
import Link from "next/link";
import { getFrameUrl } from '../utils';

export default function FrameLink({ episode, timestamp }: { episode: string, timestamp: number }) {
    return (
        <Link href={`/episode/${episode}/${timestamp}`}>
            <Image 
                src={getFrameUrl(episode, timestamp)} 
                alt={`${episode}: ${timestamp}`} 
                width={400} 
                height={240} 
                className="box-border hover:outline hover:outline-8 hover:outline-theme-black" 
            />
        </Link>
    )
}