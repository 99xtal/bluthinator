import { dehydrate, HydrationBoundary, QueryClient } from "@tanstack/react-query";
import { getFrame, getNearbyFrames } from "~/api";
import FrameEditor from "~/app/components/FrameEditor";
import FrameLink from "~/app/components/FrameLink";
import NearbySelector from "~/app/components/NearbySelector";

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