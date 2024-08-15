"use client";

import FrameLink from "./FrameLink";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getNearbyFrames } from "~/api";

type Props = {
    episode: string;
    currentTimestamp: number;
}

export default function NearbySelector({ episode, currentTimestamp }: Props) {
    const [nearbyTimestamp, setNearbyTimestamp] = useState<number>(currentTimestamp)
    const { data } = useQuery({ 
        queryKey: ['nearby', episode, nearbyTimestamp],
        queryFn: async () => getNearbyFrames(episode, nearbyTimestamp)
    })

    return (
        <div className="flex flex-row gap-4">
            {data?.map((frame) => (
                <FrameLink 
                    key={frame.id} 
                    episode={frame.episode} 
                    timestamp={frame.timestamp}
                    className={frame.timestamp !== currentTimestamp ? "filter grayscale hover:grayscale-0" : ''}
                />
            ))}
        </div>
    )
}