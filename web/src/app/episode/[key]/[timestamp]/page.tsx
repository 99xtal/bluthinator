import Image from "next/image";
import Link from "next/link";
import { defonteRegular } from "~/app/fonts";
import { getFrameUrl, msToTime } from "~/app/utils";

async function getFrame(key: string, timestamp: string) {
    const response = await fetch(`${process.env.API_HOST}/episode/${key}/${timestamp}`);
    return response.json();
}

export default async function Page({ params }: { params: { key: string, timestamp: string } }) {
    const data = await getFrame(params.key, params.timestamp);

    const caption = false;

    return (
        <div>
            <div className="flex lg:flex-row flex-col gap-8">
                <div className="relative flex flex-2">
                    <Image
                        src={getFrameUrl(data.frame.episode, data.frame.timestamp, 'large')}
                        alt={data.subtitle.text}
                        width={800}
                        height={500}
                        className="w-full h-auto outline outline-8 outline-theme-black"
                    />
                    {caption && (
                        <div className="absolute bottom-12 left-0 w-full text-center text-white p-2 ">
                            <h2 className={`${defonteRegular.className} text-2xl text-theme-yellow`}>{data.subtitle.text}</h2>
                        </div>
                    )}
                </div>
                <div className="flex flex-1 flex-col">
                    <h1 className={`${defonteRegular.className} text-3xl`}>{`"${data.episode.title}"`}</h1>
                    <div className="flex flex-row justify-between items-center">
                        <h3 className={`${defonteRegular.className} text-lg text-theme-red`}>{`Season ${data.episode.season}, Episode ${data.episode.episode_number} (${msToTime(data.frame.timestamp)})`}</h3>
                        <Link className={`${defonteRegular.className} text-md text-theme-black hover:underline`} href={`/episode/${data.frame.episode}#${data.frame.timestamp}`}>
                            <p>View Episode</p>
                        </Link>
                    </div>
                    <hr className="border-t border-gray-300 my-2" />
                    <div className="flex justify-center p-8">
                        <h3 className={`${defonteRegular.className} text-lg text-theme-black`}>{'"' + data.subtitle.text + '"'}</h3>
                    </div>
                </div>
            </div>
        </div>
    )
}