import localFont from 'next/font/local'
import Results from "./ui/Results";
import Search from "./ui/Search";

const myFont = localFont({ src: 'fonts/DeFonteReducedNormal.ttf' });

export default function Home({ searchParams}: { searchParams: { q?: string; } }) {
  const query = searchParams?.q || '';

  return (
    <main>
      <header className="flex flex-row justify-between items-center px-16 py-4" >
        <div className="transform -rotate-5">
          <h1 className={`${myFont.className} text-3xl`}>bluthinator</h1>
        </div>
        <div className="p-4 bg-theme-orange flex justify-center">
          <Search placeholder="Search for something" />
        </div>
      </header>
      <div className="flex justify-center py-4">
        <Results query={query} />
      </div>
    </main>
  );
}