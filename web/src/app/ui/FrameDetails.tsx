"use client"

import Image from "next/image"
import { getFrameUrl, msToTime } from "../utils"
import { Episode, Frame, Subtitle } from "~/types";
import { useState } from "react";
import { defonteRegular } from "../fonts";
import Link from "next/link";

export default function FrameDetails({ frame, episode, subtitle }: { frame: Frame, episode: Episode, subtitle?: Subtitle }) {
    const [isMemeMode, setMemeMode] = useState(false)
    const [caption, setCaption] = useState(subtitle?.text || '')

    const handleCancel = () => {
        setMemeMode(false)
        setCaption(subtitle?.text || '')
    }

    return (
		<div className="flex lg:flex-row flex-col gap-8">
			<div className="relative flex flex-1">
				<Image
					src={getFrameUrl(frame.episode, frame.timestamp, 'large')}
					alt={subtitle?.text || frame.episode + ' ' + frame.timestamp}
					width={640}
					height={360}
					className="w-full h-auto outline outline-8 outline-theme-black"
				/>
				{isMemeMode && (
					<div className="absolute bottom-8 left-0 w-full text-center text-white p-2 ">
						<h2
                            className={`${defonteRegular.className} text-xl text-theme-yellow`}
                            style={{ textShadow: '2px 2px 4px rgba(0, 0, 0, 0.9)' }}
                        >{caption}</h2>
					</div>
				)}
			</div>
			<div className="flex flex-1 flex-col">
				<h1 className={`${defonteRegular.className} text-3xl`}>{`"${episode.title}"`}</h1>
				<div className="flex flex-row justify-between items-center">
					<h3 className={`${defonteRegular.className} text-lg text-theme-red`}>{`Season ${episode.season}, Episode ${episode.episode_number} (${msToTime(frame.timestamp)})`}</h3>
					<Link className={`${defonteRegular.className} text-md text-theme-black hover:underline`} href={`/episode/${frame.episode}#${frame.timestamp}`}>
						<p>View Episode</p>
					</Link>
				</div>
				<hr className="border-t border-gray-300 my-2" />
				{subtitle && (
					<div className="flex-grow justify-center items-center p-8">
                        {isMemeMode ? (
                            <textarea
                                value={caption}
                                onChange={(e) => setCaption(e.target.value)}
                                className={`${defonteRegular.className} w-full p-2 text-lg text-theme-black bg-transparent border border-gray-300 outline-none resize-none`}
                                rows={3}
                            />
                            ) : (
                                <h3 className={`${defonteRegular.className} text-lg text-theme-black`}>
                                    {'"' + subtitle?.text + '"'}
                                </h3>
                            )
                        }
					</div>
				)}
                <hr className="border-t border-gray-300 my-2" />
                <div className="flex flex-row gap-2">
                    {!isMemeMode && <button className={`${defonteRegular.className} flex-1 bg-theme-red text-white p-2 rounded-md`} onClick={() => setMemeMode(true)}>
                        Make Meme
                    </button>}
                    {isMemeMode && (
                        <>
                            <button className={`${defonteRegular.className} flex-1 bg-theme-red text-white p-2 rounded-md`}>
                                Generate Meme
                            </button>
                            <button className={`${defonteRegular.className} flex-1 bg-theme-red text-white p-2 rounded-md`} onClick={handleCancel}>
                                Cancel
                            </button>
                        </>
                    )}
                    {!isMemeMode && <button className={`${defonteRegular.className} flex-1 bg-theme-red text-white p-2 rounded-md`}>
                        Make GIF
                    </button>}
                </div>
			</div>
		</div>
    )
}