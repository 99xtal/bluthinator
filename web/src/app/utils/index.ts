export const getFrameUrl = (episode: string, timestamp: number, size: 'small' | 'medium' | 'large' = 'small') => {
    return `${process.env.NEXT_PUBLIC_IMG_HOST}/frames/${episode}/${timestamp}/${size}.png`;
}