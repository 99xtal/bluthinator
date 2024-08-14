import FrameLink from "./components/FrameLink";

async function search(query: string): Promise<any> {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/search?q=${query}`);
  return response.json();
}

export default async function Home({ searchParams}: { searchParams: { q?: string; } }) {
  const query = searchParams?.q || '';
  const data = await search(query);

  return (
    <div className="flex justify-center">
      <div className="grid grid-cols-3 gap-2">
        {data.map((result: any, i: number) => (
          <div key={result.timestamp + result.episode + i} className="p-1">
            <FrameLink episode={result.episode} timestamp={result.timestamp} />
          </div>
        ))}
      </div>
    </div>
  );
}