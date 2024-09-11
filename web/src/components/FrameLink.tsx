import Image from 'next/image';
import Link from "next/link";
import { getFrameUrl } from '~/utils';

interface Props {
    episode: string;
    timestamp: number;
    size?: 'small' | 'medium' | 'large';
    className?: string;
}

export default function FrameLink({ episode, timestamp, size = 'small', className }: Props) {
    return (
        <Link href={`/episode/${episode}/${timestamp}`}>
            <Image 
                src={getFrameUrl(episode, timestamp, size)} 
                alt={`${episode}: ${timestamp}`} 
                width={400} 
                height={240}
                className={"box-border hover:outline hover:outline-4 hover:outline-theme-black bg-gray-300" + (className ? ` ${className}` : '')}
                blurDataURL={getFrameUrl(episode, timestamp, 'small')}
            />
        </Link>
    )
}