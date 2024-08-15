export interface Frame {
    id: number;
    episode: string;
    timestamp: number;
}

export interface Episode {
    id: number;
    title: string;
    season: number;
    episode_number: number;
    director: string;
}

export type EpisodeData = {
    episode_number: number;
    season: number;
    title: string;
    director: string;
    subtitles: {
        id: number;
        episode: string;
        text: string;
        start_timestamp: number;
        end_timestamp: number;
        frame_timestamp: number;
    }[]
}

export interface Subtitle {
    id: number;
    episode: string;
    text: string;
    start_timestamp: number;
    end_timestamp: number;
}