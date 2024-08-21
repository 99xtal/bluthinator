"use client"

import Image from "next/image"
import { useState } from "react";
import Link from "next/link";

import { defonteRegular } from "~/fonts";
import { getFrameUrl, msToTime } from "~/utils"
import { Episode, Frame, Subtitle } from "~/types";
import { Button, Caption, Divider, SubtitleText, TextLink, TitleText } from "~/app/elements";

export default function FrameEditor({ frame, episode, subtitle }: { frame: Frame, episode: Episode, subtitle?: Subtitle }) {
    const [isMemeMode, setMemeMode] = useState(false)
    const [caption, setCaption] = useState(subtitle?.text || '')
    
    const handleCancel = () => {
        setMemeMode(false)
        setCaption(subtitle?.text || '')
    }

    return (
		<div className="flex lg:flex-row flex-col gap-8">
			<div className="relative flex flex-1">
                <Link href={`/img/${frame.episode}/${frame.timestamp}/large.jpg`} className="block w-full">
                    <Image
                        src={getFrameUrl(frame.episode, frame.timestamp, 'large')}
                        alt={subtitle?.text || frame.episode + ' ' + frame.timestamp}
                        width={640}
                        height={360}
                        className="w-full h-auto outline outline-4 outline-theme-black"
                    />
                </Link>
				{isMemeMode && (
					<div className="absolute bottom-4 left-0 w-full text-center px-2 ">
                        <Caption>{caption}</Caption>
					</div>
				)}
			</div>
			<div className="flex flex-1 flex-col">
                <TitleText>{`"${episode.title}"`}</TitleText>
				<div className="flex flex-row justify-between items-center">
                    <SubtitleText>{`Season ${episode.season}, Episode ${episode.episode_number} (${msToTime(frame.timestamp)})`}</SubtitleText>
					<TextLink href={`/episode/${frame.episode}`}>
                        View Episode
                    </TextLink>
				</div>
                <Divider />
                <div className="flex-grow justify-center items-center p-8">
                    {isMemeMode ? (
                        <textarea
                            value={caption}
                            onChange={(e) => setCaption(e.target.value)}
                            className={`${defonteRegular.className} w-full p-2 text-lg text-theme-black bg-transparent border border-gray-300 outline-none resize-none`}
                            rows={3}
                        />
                    ) : (
                        subtitle && <h3 className={`${defonteRegular.className} text-lg text-theme-black`}>
                            {'"' + subtitle?.text + '"'}
                        </h3>
                    )
                    }
                </div>
                <Divider />
                <div className="flex flex-row gap-2">
                    {!isMemeMode && 
                        <Button onClick={() => setMemeMode(true)}>
                            Make Meme
                        </Button>
                    }
                    {isMemeMode && (
                        <>
                            <Link href={`/meme/${frame.episode}/${frame.timestamp}/${btoa(caption)}`} className={`${defonteRegular.className} flex-1 bg-theme-red text-white p-2 rounded-md flex justify-center items-center`}>
                                Generate Meme
                            </Link>
                            <Button onClick={handleCancel}>
                                Cancel
                            </Button>
                        </>
                    )}
                    {/* {!isMemeMode && 
                        <Button>
                            Make GIF
                        </Button>
                    } */}
                </div>
			</div>
		</div>
    )
}