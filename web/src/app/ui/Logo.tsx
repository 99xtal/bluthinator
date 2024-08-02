import localFont from "next/font/local";

const myFont = localFont({ src: '../fonts/DeFonteReducedNormal.ttf' });

export default function Logo() {
    return (
        <div className="transform -rotate-5">
            <h1 className={`${myFont.className} text-3xl`}>bluthinator</h1>
        </div>
    );
}