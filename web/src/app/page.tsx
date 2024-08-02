import Results from "./ui/Results";

export default function Home({ searchParams}: { searchParams: { q?: string; } }) {
  const query = searchParams?.q || '';

  return (
      <div className="flex justify-center">
        <Results query={query} />
      </div>
  );
}