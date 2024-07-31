async function search(query: string): Promise<any> {
    const response = await fetch(`${process.env.API_HOST}/search?q=${query}`);
    return response.json();
}

export default async function Results({ query }: { query: string }) {
    const data = await search(query);

    return (<pre>{JSON.stringify(data, null, 2)}</pre>)
}