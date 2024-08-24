"use client";

import { useEffect, useState } from "react";
import TitleText from "../elements/TitleText";

export default function Logo() {
    const [isLargeScreen, setIsLargeScreen] = useState(false);

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
        <div className="transform -rotate-5">
            <TitleText>{isLargeScreen ? "bluthinator" : "b"}</TitleText>
        </div>
    );
}