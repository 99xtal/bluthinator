import Link from "next/link";
import FrameLink from "~/components/FrameLink";

async function search(query: string): Promise<any> {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/search?q=${query}`);
  return response.json();
}

export default async function Home({ searchParams}: { searchParams: { q?: string; } }) {
  const query = searchParams?.q || '';
  const data = await search(query);

  return (
    <div className="flex flex-col min-h-screen">
      <div className="flex-grow flex justify-center">
        <div className="grid grid-cols-2 lg:grid-cols-3 gap-0 lg:gap-2">
          {data.map((result: any, i: number) => (
            <div key={result.timestamp + result.episode + i} className="p-1">
              <FrameLink episode={result.episode} timestamp={result.timestamp} size="medium" />
            </div>
          ))}
        </div>
      </div>
      <footer className="container mx-auto bg-theme-white py-4 flex justify-between">
        <div className="text-center text-gray-500 text-sm">
          Created By&nbsp;<a href="https://www.99xtal.com" className="underline">7\</a>
        </div>
        <Link href="/about" className="text-center text-gray-500 text-sm underline">
          About
        </Link>
      </footer>
    </div>
  );
}