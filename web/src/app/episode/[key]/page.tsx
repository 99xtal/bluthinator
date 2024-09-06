import { ffBlurProMedium } from "~/assets/fonts";

import { getEpisode } from "~/api";
import { FrameLink, ScrollToAnchor, ScrollToTopButton } from "~/components";
import { SubtitleText, TitleText } from "~/elements";
import { msToTime } from "~/utils";

interface Props {
    params: {
        key: string;
    };
}

export function generateMetadata({ params }: Props) {
    return {
        title: `Bluthinator | ${params.key}`,
    }
}

export default async function Page({ params }: Props) {
    const data = await getEpisode(params.key);

    return (
        <div>
            <ScrollToAnchor />
            <TitleText>{`Episode ${data.season}x${data.episode_number} - "${data.title}"`}</TitleText>
            <SubtitleText>{`Director: ${data.director}`}</SubtitleText>
            <div>
                {data.subtitles.map((subtitle, i) => (
                    <div key={subtitle.id} id={subtitle.frame_timestamp.toString()} className="p-4">
                        <div className="flex flex-row gap-6 items-center">
                            <div className="flex flex-1 justify-center items-end">
                                <FrameLink episode={subtitle.episode} timestamp={subtitle.frame_timestamp} size="medium" />
                            </div>
                            <div className="flex flex-1 flex-col justify-center">
                                <div className={`transform ${i % 2 === 0 ? 'rotate-2.5' : '-rotate-2.5'}`}>
                                    <h1 className={`${ffBlurProMedium.className} text-xl`}>{subtitle.text}</h1>
                                </div>
                                <div className={`transform ${i % 2 === 0 ? 'rotate-2.5' : '-rotate-2.5'}`}>
                                    <p className={`${ffBlurProMedium.className} text-theme-red`}>{msToTime(subtitle.start_timestamp)} - {msToTime(subtitle.end_timestamp)}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
            <ScrollToTopButton />
        </div>
    )
}