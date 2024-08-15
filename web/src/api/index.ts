import { EpisodeData, Frame } from "~/types";

export async function getFrame(key: string, timestamp: string) {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/episode/${key}/${timestamp}`);
    return response.json();
}
  
export async function getNearbyFrames(key: string, timestamp: number): Promise<Frame[]> {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/nearby?e=${key}&t=${timestamp}`);
    return response.json() as Promise<Frame[]>;
}

export async function getEpisode(key: string) {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/episode/${key}`);
    return response.json() as Promise<EpisodeData>;
}