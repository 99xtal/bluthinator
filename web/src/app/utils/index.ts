export const getFrameUrl = (episode: string, timestamp: number, size: 'small' | 'medium' | 'large' = 'small') => {
    return `${process.env.NEXT_PUBLIC_IMG_HOST}/frames/${episode}/${timestamp}/${size}.png`;
}

export const msToTime = (ms: number) => {
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor((totalSeconds % 3600) / 60).toString().padStart(2, '0');
    const seconds = (totalSeconds % 60).toString().padStart(2, '0');
    
    return `${minutes}:${seconds}`;
}