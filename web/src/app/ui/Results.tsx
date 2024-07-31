async function search(query: string): Promise<any> {
    const response = await fetch(`${process.env.API_HOST}/search?q=${query}`);
    return response.json();
}

export default async function Results({ query }: { query: string }) {
    const data = await search(query);

    const getImgUrl = (result: any) => {
        return `${process.env.NEXT_PUBLIC_IMG_HOST}/frames/${result.episode}/${result.timestamp}/small.png`;
    }

    return (
        <div className="grid grid-cols-3 gap-4">
            {data.map((result: any) => (
                <img key={result.frame} src={getImgUrl(result)} />
            ))}
        </div>
    );
}