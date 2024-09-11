"use client"

import Image from "next/image"
import { useRouter } from "next/navigation"
import { useEffect, useRef, useState } from "react";
import Link from "next/link";

import { ffBlurProMedium } from "~/assets/fonts";
import { getFrameUrl, msToTime } from "~/utils"
import { Episode, Frame, Subtitle } from "~/types";
import { Button, Caption, Divider, SubtitleText, TextLink, TitleText } from "~/elements";
import { logEvent } from "~/utils/firebase";

const CAPTION_FONT_RATIO = 0.08;
const CAPTION_BOTTOM_RATION = 0.045;

export default function FrameEditor({ frame, episode, subtitle }: { frame: Frame, episode: Episode, subtitle?: Subtitle }) {
    const router = useRouter();
    const [isMemeMode, setMemeMode] = useState(false)
    const [caption, setCaption] = useState(subtitle?.text || '')
    const [imageSize, setImageSize] = useState({ width: 0, height: 0 })
    const captionTooLong = caption.length > 500
    const imageRef = useRef<HTMLImageElement | null>(null);

    useEffect(() => {
        const onResize = () => {
            if (imageRef.current) {
                const { width, height } = imageRef.current.getBoundingClientRect()
                setImageSize({ width, height })
            }
        }

        onResize();
        window.addEventListener('resize', onResize)

        return () => {
            window.removeEventListener('resize', onResize)
        }
    }, [])

    const handleMakeMeme = () => {
        setMemeMode(true)
        logEvent('edit_meme', { episode: frame.episode, timestamp: frame.timestamp })
    }
    
    const handleCancel = () => {
        setMemeMode(false)
        setCaption(subtitle?.text || '')
        logEvent('cancel_meme', { episode: frame.episode, timestamp: frame.timestamp })
    }

    const handleGenerateMeme = () => {
        logEvent('generate_meme', { episode: frame.episode, timestamp: frame.timestamp, caption })
        router.push(`/meme/${frame.episode}/${frame.timestamp}/${btoa(caption)}`)
    }

    return (
		<div className="flex lg:flex-row flex-col gap-8">
			<div className="relative flex flex-1">
                <Link href={`/img/${frame.episode}/${frame.timestamp}/large.jpg`} className="block w-full">
                    <Image
                        ref={imageRef}
                        src={getFrameUrl(frame.episode, frame.timestamp, 'large')}
                        alt={subtitle?.text || frame.episode + ' ' + frame.timestamp}
                        width={640}
                        height={360}
                        className="w-full h-auto outline outline-4 outline-theme-black"
                    />
                </Link>
				{isMemeMode && (
					<div style={{ bottom: `${Math.floor(imageSize.height * CAPTION_BOTTOM_RATION)}px`}} className="absolute left-0 w-full text-center px-2 ">
                        <Caption fontSize={Math.floor(imageSize.height * CAPTION_FONT_RATIO)}>{caption}</Caption>
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
                <div className="flex flex-col flex-grow justify-center items-center gap-2 p-8">
                    {isMemeMode ? (
                        <>
                            <textarea
                                value={caption}
                                onChange={(e) => setCaption(e.target.value)}
                                className={`${ffBlurProMedium.className} w-full p-2 text-lg text-theme-black bg-transparent border ${captionTooLong ? "border-theme-red" : "border-gray-300"} outline-none resize-none`}
                                rows={3}
                            />
                            <p className={`self-end text-xs ${captionTooLong ? 'text-theme-red' : 'text-gray-500' }`}>{`(${caption.length}/500)`}</p>
                        </>
                    ) : (
                        subtitle && <h3 className={`${ffBlurProMedium.className} text-lg text-theme-black`}>
                            {'"' + subtitle?.text + '"'}
                        </h3>
                    )
                    }
                </div>
                <Divider />
                <div className="flex flex-row gap-2">
                    {!isMemeMode && 
                        <Button onClick={handleMakeMeme} className="flex-1">
                            Make Meme
                        </Button>
                    }
                    {isMemeMode && (
                        <>
                            <Button onClick={handleCancel} variant='secondary' className="flex-1">
                                Cancel
                            </Button>
                            <Button 
                                onClick={handleGenerateMeme} 
                                disabled={captionTooLong}
                                className={`${ffBlurProMedium.className} flex-1 bg-theme-red text-white p-2 rounded-md flex justify-center items-center`}
                            >
                                Generate
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