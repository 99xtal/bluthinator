import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query";
import { Metadata } from "next";
import { getFrame, getNearbyFrames } from "~/api";
import { FrameEditor, NearbySelector } from "~/components";

interface Props {
  params: { key: string, timestamp: string }
}

export async function generateMetadata(
  { params }: Props,
): Promise<Metadata> {
  return {
    openGraph: {
      title: "Bluthinator",
      description: "An Arrested Development search engine and meme generator",
      images: [
        {
          url: `${process.env.IMG_HOST}/frames/${params.key}/${params.timestamp}/large.jpg`,
          width: 720,
          height: 405,
          alt: "Bluthinator image",
        },
      ]
    }
  }
}

export default async function Page({ params }: { params: { key: string, timestamp: string } }) {
  const data = await getFrame(params.key, params.timestamp);
  const queryClient = new QueryClient()
  await queryClient.prefetchQuery({ 
    queryKey: ['nearby', params.key, parseInt(params.timestamp)],
    queryFn: async () => getNearbyFrames(params.key, parseInt(params.timestamp))
  })
  const dehydratedState = dehydrate(queryClient)

  return (
    <div className="flex flex-col gap-16">
        <FrameEditor frame={data.frame} episode={data.episode} subtitle={data.subtitle} />
        <HydrationBoundary state={dehydratedState}>
            <NearbySelector episode={params.key} currentTimestamp={data.frame.timestamp} />
        </HydrationBoundary>
    </div>
  )
}