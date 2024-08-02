export default function Page({ params }: { params: { key: string } }) {
    return <div>My Post: {params.key}</div>
}