'use client';

import { useEffect } from 'react';

export default function ScrollToAnchor() {
    useEffect(() => {
        const hash = window.location.hash;
        if (hash) {
            console.log(hash, hash.substring(1));
            const element = document.getElementById(hash.substring(1));
            if (element) {
                element.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    }, []);

    return null;
};
