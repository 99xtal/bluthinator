import type { Metadata } from "next";
import { Inter } from "next/font/google";
import Link from "next/link";
import "./globals.css";
import Logo from "./ui/Logo";
import Search from "./ui/Search";
import { Suspense } from "react";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Bluthinator",
  description: "An Arrested Development search engine",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${inter.className} bg-theme-white`}>
        <header className="sticky top-0 z-14 flex flex-row justify-center items-center pl-8 py-4 gap-64 bg-theme-white" >
          <Link href="/">
            <Logo />
          </Link>
          <div className="px-8 py-4 bg-theme-orange flex flex-grow">
            <Suspense>
              <Search 
                placeholder="Search for something" 
                className="w-full border border-gray-300 rounded-full p-2 focus:outline-none focus:ring-2 focus:ring-gray-400 hover:border-gray-400"
              />
            </Suspense>
          </div>
        </header>
        <main className="container mx-auto py-4">
          {children}
        </main>
      </body>
    </html>
  );
}
