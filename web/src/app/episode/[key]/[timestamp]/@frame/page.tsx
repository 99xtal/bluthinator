import FrameEditor from "~/app/components/FrameEditor";

async function getFrame(key: string, timestamp: string) {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_HOST}/episode/${key}/${timestamp}`);
  return response.json();
}

export default async function Page({ params }: { params: { key: string, timestamp: string } }) {
  const data = await getFrame(params.key, params.timestamp);

  return (
    <FrameEditor frame={data.frame} episode={data.episode} subtitle={data.subtitle} />
  )
}