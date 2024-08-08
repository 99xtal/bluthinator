import FrameDetails from "~/app/ui/FrameDetails";

async function getFrame(key: string, timestamp: string) {
    const response = await fetch(`${process.env.API_HOST}/episode/${key}/${timestamp}`);
    return response.json();
}

export default async function Page({ params }: { params: { key: string, timestamp: string } }) {
    const data = await getFrame(params.key, params.timestamp);

    return (
		<FrameDetails frame={data.frame} episode={data.episode} subtitle={data.subtitle} />
    )
}