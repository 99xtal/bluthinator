"use client";

import FrameLink from "./FrameLink";
import { useState } from "react";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import { getNearbyFrames } from "~/api";

type Props = {
    episode: string;
    currentTimestamp: number;
}

export default function NearbySelector({ episode, currentTimestamp }: Props) {
    const [nearbyTimestamp, setNearbyTimestamp] = useState<number>(currentTimestamp)
    const { data } = useQuery({ 
        queryKey: ['nearby', episode, nearbyTimestamp],
        queryFn: async () => getNearbyFrames(episode, nearbyTimestamp),
        placeholderData: keepPreviousData
    })

    if (!data) return null

    return (
        <div className="flex flex-row gap-4">
            <button 
                onClick={() => setNearbyTimestamp(data[0].timestamp)}
                className="p-2 md:p-4 text-3xl"
            >
                &lt;
            </button>
            <div className="grid grid-cols-2 gap-2 md:grid-cols-7">
                {data.map((frame) => (
                    <FrameLink 
                        key={frame.id} 
                        episode={frame.episode} 
                        timestamp={frame.timestamp}
                        className={frame.timestamp !== currentTimestamp ? "filter grayscale hover:grayscale-0" : ''}
                    />
                ))}
            </div>
            <button 
                onClick={() => setNearbyTimestamp(data[data.length - 1].timestamp)}
                className="p-2 md:p-4 text-3xl"
            >
                &gt;
            </button>
        </div>
    )
}