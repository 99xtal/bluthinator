"use client";

import Image from "next/image";
import { useEffect, useState } from "react";
import TitleText from "../elements/TitleText";
import circleSvg from "~/assets/svg/circle.svg";

export default function Logo() {
    const [isLargeScreen, setIsLargeScreen] = useState<boolean>();

    useEffect(() => {
        const mediaQuery = window.matchMedia('(min-width: 640px)');
        const handleMediaQueryChange = (event: MediaQueryListEvent) => {
            setIsLargeScreen(event.matches);
        };

        setIsLargeScreen(mediaQuery.matches);

        mediaQuery.addEventListener('change', handleMediaQueryChange);

        return () => {
            mediaQuery.removeEventListener('change', handleMediaQueryChange);
        };
    }, []);

    return (
        <div className="transform -rotate-5 relative inline-block px-8 py-4">
            <Image 
                src={circleSvg}
                alt=""
                className={`absolute transform ${isLargeScreen ? "scale-75 scale-x-90 -top-7 -left-0" : "top-2 -left-0 scale-125 scale-x-75"} -z-10`}
            />
            <TitleText className="relative">{isLargeScreen ? "bluthinator" : "b"}</TitleText>
        </div>
    );
}