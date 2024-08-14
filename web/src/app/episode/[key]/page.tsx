import { defonteRegular } from "~/fonts";

import FrameLink from "~/app/components/FrameLink";
import ScrollToAnchor from "~/app/components/ScrollToAnchor";
import { msToTime } from "~/utils";

type EpisodeData = {
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

async function getEpisode(key: string) {
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/episode/${key}`);
    return response.json() as Promise<EpisodeData>;
}

export default async function Page({ params }: { params: { key: string } }) {
    const data = await getEpisode(params.key);

    return (
        <div>
            <ScrollToAnchor />
            <h1 className={`${defonteRegular.className} text-2xl`}>Episode {data.season}x{data.episode_number} - {`"${data.title}"`}</h1>
            <h3 className={`${defonteRegular.className} text-xl text-theme-red`}>Director: {data.director}</h3>
            <div>
                {data.subtitles.map((subtitle, i) => (
                    <div key={subtitle.id} id={subtitle.frame_timestamp.toString()} className="p-4">
                        <div className="flex flex-row gap-2">
                            <div className="flex flex-1 justify-center items-end">
                                <FrameLink episode={subtitle.episode} timestamp={subtitle.frame_timestamp} />
                            </div>
                            <div className="flex flex-1 flex-col justify-center">
                                <div className={`transform ${i % 2 === 0 ? 'rotate-2.5' : '-rotate-2.5'}`}>
                                    <h1 className={`${defonteRegular.className} text-xl`}>{subtitle.text}</h1>
                                </div>
                                <div className={`transform ${i % 2 === 0 ? 'rotate-2.5' : '-rotate-2.5'}`}>
                                    <p className={`${defonteRegular.className} text-theme-red`}>{msToTime(subtitle.start_timestamp)} - {msToTime(subtitle.end_timestamp)}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    )
}