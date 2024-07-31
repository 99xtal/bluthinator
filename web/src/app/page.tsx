import Results from "./ui/Results";
import Search from "./ui/Search";

export default function Home({ searchParams}: { searchParams: { q?: string; } }) {
  const query = searchParams?.q || '';

  return (
    <main>
      <Search placeholder="Search for something" />
      <Results query={query} />
    </main>
  );
}