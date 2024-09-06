import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Link from "next/link";
import "./globals.css";
import Logo from "~/components/Logo";
import Search from "~/components/Search";
import { Suspense } from "react";
import Providers from "./providers";
import Script from "next/script";
import { Button } from "../elements";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Bluthinator",
  description: "An Arrested Development search engine",
};

const references = [
  "There's always money in the banana stand.",
  "I've made a huge mistake.",
  "I'm a monster!",
  "Her?",
  "Hey brother",
  "I'm afraid I just blue myself.",
  "That was a freebie.",
  "Marry me!",
  "No touching!",
]

function genSearchPlaceholder() {
  return `"${references[Math.floor(Math.random() * references.length)]}"`;
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className} bg-theme-white`}>
        <Script id="facebook-sdk-init">
          {`window.fbAsyncInit = function() {
            FB.init({
              appId            : '${process.env.NEXT_PUBLIC_FACEBOOK_APP_ID}',
              xfbml            : true,
              version          : 'v20.0'
            });
          };`}
        </Script>
        <Script async defer crossOrigin="anonymous" src="https://connect.facebook.net/en_US/sdk.js"></Script>
        <header className="z-50 sticky top-0 flex flex-row justify-center items-center bg-white border-black border-b-4" >
            <Link href="/">
                <Logo />
            </Link>
          <div className="px-4 md:px-8 py-4 bg-theme-orange flex flex-grow gap-2">
            <Suspense>
              <Search 
                placeholder={genSearchPlaceholder()} 
                className="w-full border border-gray-300 rounded-full p-2 focus:outline-none focus:ring-2 focus:ring-gray-400 hover:border-gray-400"
              />
            </Suspense>
            <Button>
              <a href="/random">Random</a>
            </Button>
          </div>
        </header>
        <main className="container mx-auto p-4 lg:p-8">
          <Providers>
            {children}
          </Providers>
        </main>
      </body>
    </html>
  );
}
