import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Link from "next/link";
import "./globals.css";
import Logo from "./components/Logo";
import Search from "./components/Search";
import { Suspense } from "react";
import Providers from "./providers";

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
  return `Search for quotes (like "${references[Math.floor(Math.random() * references.length)]}")`;
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className} bg-theme-white`}>
        <header className="sticky top-0 flex flex-row justify-center items-center bg-white border-b border-black border-b-4" >
          <div className="hidden md:block px-16">
            <Link href="/">
                <Logo />
            </Link>
          </div>
          <div className="px-8 py-4 bg-theme-orange flex flex-grow">
            <Suspense>
              <Search 
                placeholder={genSearchPlaceholder()} 
                className="w-full border border-gray-300 rounded-full p-2 focus:outline-none focus:ring-2 focus:ring-gray-400 hover:border-gray-400"
              />
            </Suspense>
          </div>
        </header>
        <main className="container mx-auto p-4 lg:pd-8">
          <Providers>
            {children}
          </Providers>
        </main>
      </body>
    </html>
  );
}
