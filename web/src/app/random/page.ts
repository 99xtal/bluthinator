import { redirect } from "next/navigation";
import { getRandomFrame } from "~/api";

export default async function Page() {
    const randomFrame = await getRandomFrame();

    if (randomFrame) {
        redirect(`/episode/${randomFrame.episode}/${randomFrame.timestamp}`);
    }

    return null;
}