import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import Logo from "./ui/Logo";
import Search from "./ui/Search";

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
        <header className="sticky top-0 flex flex-row justify-between items-center px-16 py-4 bg-theme-white" >
          <Logo />
          <div className="p-4 bg-theme-orange flex justify-center">
            <Search placeholder="Search for something" />
          </div>
        </header>
        <main className="container mx-auto py-4">
          {children}
        </main>
      </body>
    </html>
  );
}
