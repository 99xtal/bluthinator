import Logo from "./ui/Logo";
import Results from "./ui/Results";
import Search from "./ui/Search";

export default function Home({ searchParams}: { searchParams: { q?: string; } }) {
  const query = searchParams?.q || '';

  return (
    <main>
      <header className="flex flex-row justify-between items-center px-16 py-4" >
        <Logo />
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