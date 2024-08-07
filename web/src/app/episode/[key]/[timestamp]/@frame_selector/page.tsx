export default async function Page({ params }: { params: { key: string, timestamp: string } }) {
    return (
        <div>
            <h1>{`Selector ${params.key} ${params.timestamp}`}</h1>
        </div>
    )
}