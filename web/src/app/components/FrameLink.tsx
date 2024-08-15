import Image from 'next/image';
import Link from "next/link";
import { getFrameUrl } from '~/utils';

interface Props {
    episode: string;
    timestamp: number;
    className?: string;
}

export default function FrameLink({ episode, timestamp, className }: Props) {
    return (
        <Link href={`/episode/${episode}/${timestamp}`}>
            <Image 
                src={getFrameUrl(episode, timestamp)} 
                alt={`${episode}: ${timestamp}`} 
                width={400} 
                height={240}
                className={"box-border hover:outline hover:outline-4 hover:outline-theme-black" + (className ? ` ${className}` : '')}
            />
        </Link>
    )
}