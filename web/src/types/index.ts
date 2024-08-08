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

export interface Subtitle {
    id: number;
    episode: string;
    text: string;
    start_timestamp: number;
    end_timestamp: number;
}