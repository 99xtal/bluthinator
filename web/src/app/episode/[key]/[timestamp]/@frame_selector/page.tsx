import FrameLink from "~/app/ui/FrameLink";
import { Frame } from "~/types";

async function getNearbyFrames(key: string, timestamp: string): Promise<Frame[]> {
    const response = await fetch(`${process.env.API_HOST}/nearby?e=${key}&t=${timestamp}`);
    return response.json() as Promise<Frame[]>;
}

export default async function Page({ params }: { params: { key: string, timestamp: string } }) {
    const frames = await getNearbyFrames(params.key, params.timestamp);
    
    return (
        <div className="flex flex-row gap-4">
            {frames.map((frame) => (
                <FrameLink 
                    key={frame.id} 
                    episode={frame.episode} 
                    timestamp={frame.timestamp}
                    className={frame.timestamp !== parseInt(params.timestamp) ? "filter grayscale hover:grayscale-0" : ''}
                />
            ))}
        </div>
    )
}